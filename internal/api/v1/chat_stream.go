package v1

import (
	"errors"
	"gorm.io/gorm"
	"io"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
)

const HeaderUserIp = "X-User-IP"

// Stream godoc
// @Summary      流式传输文本
// @Description  将会通过 SSE 的方式来流式传输内容，不建议使用本文档生成的代码来获取，第三方库有更好的解决方案。
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     none
// @Param 	     X-User-IP  header  string  false  "指定聊天中的用户 IP 地址，不指定则自动获取。此 IP 地址只会增加至 Prompt 中，如果不希望增加，请关闭系统自带 Prompt 选项"
// @Param        stream_id  path  string  true  "Chat stream id"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/stream/{stream_id} [get]
func (u *ChatController) Stream(c *gin.Context) {
	var response = schema.NewResponse(c)
	// 检查 stream id 是否存在
	streamIdStr := c.Param("stream_id")
	streamIdCacheKey := u.getCacheKey("stream:" + streamIdStr)
	// 检查缓存是否存在
	i, err := u.redis.Exists(c, streamIdCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if i == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatStreamNotFound).Send()
		return
	}

	// 获取 chat id
	chatIdStr, err := u.redis.Get(c, streamIdCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var chatIdStreamKey = u.getChatIdStreamingKey(schema.EntityId(chatId))
	// 检查状态是否是回复中
	isStreaming := u.isStreaming(c, schema.EntityId(chatId))
	if isStreaming {
		response.Status(http.StatusConflict).Error(consts.ErrChatStreamingPleaseWait).Send()
		return
	} else {
		// 不在回复中，则设置缓存键，防止再次请求
		cmd := u.redis.Set(c, chatIdStreamKey, 1, consts.ChatStreamExpire)
		if cmd.Err() != nil {
			response.Status(http.StatusInternalServerError).Error(cmd.Err()).Send()
			return
		}
	}

	chatEntity, err := u.chatService.GetChat(c, schema.EntityId(chatId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if chatEntity.Id == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatNotFound).Send()
		return
	}

	var assistantEntity *entity.Assistant
	var tools []llms.Tool

	// 获取上一条消息，拿到指定的 assistant id
	lastChatMessage, err := u.cm.GetLatestMessage(c, chatEntity)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if lastChatMessage != nil && lastChatMessage.AssistantId != nil && (lastChatMessage.Role == schema.RoleHuman || lastChatMessage.Role == schema.RoleHideHuman) {
		assistantEntity, err = u.assistantService.GetAssistant(c, *lastChatMessage.AssistantId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		// 这里不用判断是不是用户的，因为添加消息时已经判断了
	}

	if assistantEntity == nil && chatEntity.AssistantId != nil {
		assistantEntity, err = u.assistantService.GetAssistant(c, *chatEntity.AssistantId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
		// 获取 assistant 绑定的 tools
		tools, err = u.assistantService.ToLLMTool(c, assistantEntity)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

	}

	// 提取 history
	histories, err := u.cm.GetChatMessageWithHide(c, chatEntity)
	var llmResponseChan = make(chan *schema.AssistantResponse, 1)

	streamUserCacheKey := u.getCacheKey("stream:" + streamIdStr + ":user")

	// 检查缓存是否存在
	i, err = u.redis.Exists(c, streamUserCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if i == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatStreamNotFound).Send()
		return
	}

	userCmd, err := u.redis.Get(c, streamUserCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	user := &schema.UserPublicInfo{}
	err = sonic.Unmarshal([]byte(userCmd), user)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	// MessageList
	var messageList = make([]entity.ChatMessage, 0)

	prompt, err := u.getPrompt(c, assistantEntity, user, chatEntity.Owner)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var llmChat = &schema.LLMChat{
		ResponseChan:   llmResponseChan,
		SystemPrompt:   prompt,
		UserPublicInfo: user,
		MaxTokens:      u.config.LLM.MaxTokens,
		Tools:          tools,
		Chat: &schema.ChatPublicModel{
			Name:        chatEntity.Name,
			AssistantId: chatEntity.AssistantId,
			ExpiredAt:   chatEntity.ExpiredAt,
			Owner:       chatEntity.Owner,
			GuestId:     chatEntity.GuestId,
		},
	}

	// SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	go func() {
		err = u.llmService.StreamChat(c, llmChat, histories)
		if err != nil {
			u.logger.Sugar.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()

	var llmFullMessage = ""
	var tokenUsage = &schema.TokenUsage{
		PromptTokens:     0,
		CompletionTokens: 0,
		TotalTokens:      0,
	}
	c.Stream(func(w io.Writer) bool {
		// Emit Server Sent Events compatible
		msg, ok := <-llmResponseChan
		if !ok {
			return false
		}

		if msg == nil {
			return true
		}

		j, err := sonic.MarshalString(msg)
		if err != nil {
			u.logger.Sugar.Error(err)
			return false
		}

		c.SSEvent(eventName, j)

		c.Writer.Flush()

		switch msg.State {
		case schema.StateChunk:
			llmFullMessage += msg.ChunkMessage.Content
			return true
		case schema.StateToolSuccess:
			return true
		case schema.StateToolCalling:
			var cm = entity.ChatMessage{
				Role:     schema.RoleToolCall,
				Content:  "",
				ChatId:   chatEntity.Id,
				ToolCall: (*schema.ToolCall)(msg.Internal.ToolCall),
			}

			messageList = append(messageList, cm)

			return true
		case schema.StateToolResponse:
			var cm = entity.ChatMessage{
				Role:     schema.RoleTool,
				Content:  msg.ToolResponseMessage.Content,
				ChatId:   chatEntity.Id,
				ToolCall: (*schema.ToolCall)(msg.Internal.ToolCall),
			}

			messageList = append(messageList, cm)

			// 如果有新增
			if msg.ToolResponseMessage.Append {
				cm = entity.ChatMessage{
					Role:    msg.ToolResponseMessage.Role,
					Content: msg.ToolResponseMessage.Text,
					ChatId:  chatEntity.Id,
				}

				messageList = append(messageList, cm)
			}

			return true
		case schema.StateDone:
			return true
		case schema.StateFailed:
			// 这样会发生悬垂
			// 所以要添加个新的消息
			var cm = entity.ChatMessage{
				Role:     schema.RoleTool,
				Content:  msg.ToolResponseMessage.Content,
				ChatId:   chatEntity.Id,
				ToolCall: (*schema.ToolCall)(msg.Internal.ToolCall),
			}

			messageList = append(messageList, cm)

			return false
		case schema.StateFinished:
			tokenUsage = msg.TokenUsage
			return false
		case schema.StateToolFailed:
			return false
		default:
			return true
		}
	})

	// 发送 [DONE]
	c.SSEvent(eventName, eventDone)

	// close sse stream
	c.SSEvent("close", "")
	c.Writer.Flush()

	// 移除缓存
	u.redis.Del(c, streamIdCacheKey)
	u.redis.Del(c, u.getCacheKey("entity:"+chatIdStr))
	u.redis.Del(c, u.getCacheKey("stream:"+streamIdStr+":user"))
	u.redis.Del(c, chatIdStreamKey)

	if llmFullMessage != "" {
		// 添加到消息 entity.ChatMessage
		newMessage := &entity.ChatMessage{
			Role:             schema.RoleAssistant,
			Content:          llmFullMessage,
			ChatId:           chatEntity.Id,
			PromptTokens:     tokenUsage.PromptTokens,
			CompletionTokens: tokenUsage.CompletionTokens,
			TotalTokens:      tokenUsage.TotalTokens,
		}
		messageList = append(messageList, *newMessage)
	}

	// 添加到数据库
	for _, message := range messageList {
		// 如果 assistant 不为空，则为接下来的每个消息附上当前回复的 assistant id
		if assistantEntity != nil {
			message.AssistantId = &assistantEntity.Id
		}

		err = u.cm.CreateChatMessage(c, &message)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	c.Status(http.StatusOK)
}

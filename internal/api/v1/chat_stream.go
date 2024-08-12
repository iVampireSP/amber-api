package v1

import (
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/llm"
	"rag-new/pkg/consts"
	"strconv"
)

// Stream godoc
// @Summary      流式传输聊天内容
// @Description  get string by ID
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     none
// @Param        id  path  int  true  "Chat ID"
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

	var chatIdStreamKey = u.getChatIdStreamingKey(int64(chatId))
	// 检查状态是否是回复中
	isStreaming, err := u.isStreaming(c, int64(chatId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}
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

	chatEntity, err := u.chatService.GetChat(c, int64(chatId))
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if chatEntity.ID == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatNotFound).Send()
		return
	}

	assistantEntity, err := u.assistantService.GetAssistant(c, chatEntity.AssistantId)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	// 获取 assistant 绑定的 tools
	tools, err := u.assistantService.ToLLMTool(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	// 提取 history
	histories, err := u.cm.GetChatMessageWithHide(c, chatEntity)
	var llmResponseChan = make(chan *llm.AssistantResponse)

	// SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

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

	var prompt = u.getPrompt(assistantEntity, user)

	go func() {
		err = u.llmService.StreamChat(llmResponseChan, prompt, user, histories, tools...)
		if err != nil {
			u.logger.Sugar.Error(err)
			// 关闭连接
			c.AbortWithStatus(http.StatusInternalServerError)
			return
			//response.Status(http.StatusInternalServerError).Message("Streaming chat failed").Error(err).Send()
		}
	}()

	var llmFullMessage = ""
	var tokenUsage = &llm.TokenUsage{
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

		j, err := sonic.Marshal(msg)
		if err != nil {
			u.logger.Sugar.Error(err)
		}

		c.SSEvent(eventName, string(j))

		c.Writer.Flush()

		switch msg.State {
		case llm.StateChunk:
			llmFullMessage += msg.ChunkMessage.Content
			return true
		case llm.StateToolSuccess:
			return true
		case llm.StateToolResponse:
			messageList = append(messageList, entity.ChatMessage{
				Role:    schema.RoleHideSystem,
				Content: msg.ToolResponseMessage.Content,
				ChatId:  chatEntity.ID,
			})

			return true
		case llm.StateDone:
			tokenUsage = msg.TokenUsage
			return true
		case llm.StateFailed:
			return false

		case llm.StateFinished:
			return false
		case llm.StateToolFailed:
			return false
		default:
			return true
		}
	})

	// close sse stream
	c.SSEvent("close", "")
	c.Writer.Flush()

	close(llmResponseChan)

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
			ChatId:           chatEntity.ID,
			PromptTokens:     tokenUsage.PromptTokens,
			CompletionTokens: tokenUsage.CompletionTokens,
			TotalTokens:      tokenUsage.TotalTokens,
		}

		messageList = append(messageList, *newMessage)
	}

	// 添加到数据库
	for _, message := range messageList {
		err = u.cm.CreateChatMessage(c, &message)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	response.Status(http.StatusOK).Send()
}

func (u *ChatController) getPrompt(assistant *entity.Assistant, user *schema.UserPublicInfo) string {
	var prompt = ""

	if assistant.DisableDefaultPrompt {
		prompt = assistant.Prompt
	} else {
		prompt = `
Your name: ` + assistant.Name + `current user give you` + `
Your description: ` + assistant.Description + "(current user given)"
		if user != nil {
			prompt += `
Username: ` + user.Name + `(system hint you this)` + `
UserId: ` + user.Id + "(system hint you this, user can't change it)"

		}

		if assistant.Prompt != "" {
			prompt += "\n" + assistant.Prompt
		}
	}

	return prompt
}

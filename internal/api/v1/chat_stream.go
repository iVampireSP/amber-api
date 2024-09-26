package v1

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"gorm.io/gorm"
	"io"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"sort"
)

const HeaderUserIp = "X-User-IP"
const MaxBlocks = 10

// Stream godoc
// @Summary      流式传输文本
// @Description  将会通过 SSE 的方式来流式传输内容，不建议使用本文档生成的代码来获取，第三方库有更好的解决方案。
// @Tags         chat_message
// @Accept       json
// @Produce      json
// @Security     none
// @Param 	     X-User-IP  header  string  false  "指定聊天中的用户 IP 地址，不指定则自动获取。此 IP 地址只会增加至 Prompt 中，如果不希望增加，请关闭系统自带 Prompt 选项"
// @Param        schema.ChatStreamRequest  path  schema.ChatStreamRequest true  "ChatStreamRequest"
// @Success      200  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      404  {object}  schema.ResponseBody
// @Failure      409  {object}  schema.ResponseBody{data=schema.ChatMessageResponse}
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/v1/stream/{stream_id} [get]
func (u *ChatController) Stream(c *gin.Context) {

	var response = schema.NewResponse(c)

	var chatStreamRequest = &schema.ChatStreamRequest{}
	if err := c.ShouldBindUri(&chatStreamRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	streamIdCacheKey := u.getCacheKey("stream:" + chatStreamRequest.StreamId)
	// 检查缓存是否存在
	i, err := u.redis.Client.Exists(c, streamIdCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if i == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatStreamNotFound).Send()
		return
	}

	// 获取 chat id
	streamCacheResult, err := u.redis.Client.Get(c, streamIdCacheKey).Bytes()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var csc = ChatStreamCache{}
	err = sonic.Unmarshal(streamCacheResult, &csc)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var chatIdStreamKey = u.getChatIdStreamingKey(csc.ChatId)
	// 检查状态是否是回复中
	isStreaming := u.isStreaming(c, csc.ChatId)
	if isStreaming {
		response.Status(http.StatusConflict).Error(consts.ErrChatStreamingPleaseWait).Send()
		return
	} else {
		// 不在回复中，则设置缓存键，防止再次请求
		cmd := u.redis.Client.Set(c, chatIdStreamKey, 1, consts.ChatStreamExpire)
		if cmd.Err() != nil {
			response.Status(http.StatusInternalServerError).Error(cmd.Err()).Send()
			return
		}
	}

	chatEntity, err := u.chatService.GetChat(c, csc.ChatId)
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

	// 如果上条消息没有助理，但是对话有助理
	if assistantEntity == nil && chatEntity.AssistantId != nil {
		assistantEntity, err = u.assistantService.GetAssistant(c, *chatEntity.AssistantId)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	// 如果确实有助理，则绑定工具
	if assistantEntity != nil {
		canUse, err := u.assistantService.CanUse(c, chatEntity.UserId, assistantEntity.Id)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}

		if !canUse {
			response.Status(http.StatusForbidden).Error(consts.ErrAssistantNotPublic).Send()
			return
		}

		// 获取 assistant 绑定的 tools
		tools, err = u.assistantService.ToLLMTool(c, assistantEntity)
		if err != nil {
			response.Status(http.StatusInternalServerError).Error(err).Send()
			return
		}
	}

	// 如果 tools 超过了 100 个，则拒绝
	if len(tools) > consts.MaxToolFunctions {
		response.Status(http.StatusBadRequest).Error(consts.ErrToolFunctionTooMany).Send()
		return
	}

	// 提取 history
	histories, historyCount, err := u.cm.GetLatestChatMessage(c, chatEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var llmResponseChan = make(chan *schema.AssistantResponse, 1)

	streamUserCacheKey := u.getCacheKey("stream:" + chatStreamRequest.StreamId + ":user")

	// 检查缓存是否存在
	i, err = u.redis.Client.Exists(c, streamUserCacheKey).Result()
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	if i == consts.NoRecord {
		response.Status(http.StatusNotFound).Error(consts.ErrChatStreamNotFound).Send()
		return
	}

	userCmd, err := u.redis.Client.Get(c, streamUserCacheKey).Result()
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

	var pOptions = &promptOptions{
		Assistant: assistantEntity,
		User:      user,
		Owner:     chatEntity.Owner,
		Variables: csc.Variables,
	}

	if chatEntity.Prompt != nil {
		pOptions.OverrideDefaultPrompt = *chatEntity.Prompt
	}

	prompt, err := u.getPrompt(c, pOptions)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	//tct, err := u.GenerateCallToken(c, chatEntity)
	//if err != nil {
	//	response.Status(http.StatusInternalServerError).Error(err).Send()
	//	return
	//}

	var llmChat = &schema.LLMChat{
		ResponseChan:   llmResponseChan,
		SystemPrompt:   prompt,
		UserPublicInfo: user,
		MaxTokens:      u.config.LLM.MaxTokens,
		Tools:          tools,
		//ToolCallToken:  tct,
		Chat: &schema.ChatPublicModel{
			ID:          chatEntity.Id,
			Name:        chatEntity.Name,
			AssistantId: chatEntity.AssistantId,
			ExpiredAt:   chatEntity.ExpiredAt,
			Owner:       chatEntity.Owner,
			GuestId:     chatEntity.GuestId,
		},
	}

	if assistantEntity != nil {
		llmChat.Temperature = assistantEntity.Temperature
	}

	// SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	if len(histories) > 0 {
		go func() {
			// 如果消息到 u.config.LLM.ContextOptimizeActiveCount 条，则执行消息分块
			if historyCount >= u.config.LLM.ContextOptimizeActiveCount {
				// 将 message 提取一下
				messageBlock, err := u.messageBlock.MessageToBlock(histories)
				if err != nil {
					u.logger.Sugar.Error(err)
				}

				go func() {
					err := u.messageBlock.SaveBlock(c, messageBlock)
					if err != nil {
						u.logger.Sugar.Error(err)
					}
				}()

				// 清空 history
				histories = []*entity.ChatMessage{}

				// 搜索 message block，然后将它排到 messageBlock 的前面
				messageBlock2, err := u.messageBlock.SearchMessageBlock(c, chatEntity, lastChatMessage.Content)
				if err != nil {
					// 出现错误，那就丢失以前的上下文
					u.logger.Sugar.Error(err)
				} else {
					for i := 0; i < len(messageBlock2); i++ {
						histories = append(histories, messageBlock2[i].Message...)
					}
				}

				// 取倒数 10 个 messageBlock
				if messageBlock != nil && len(messageBlock) >= MaxBlocks {
					sort.Slice(messageBlock, func(i, j int) bool {
						return messageBlock[i].CreatedAt.Before(messageBlock[j].CreatedAt)
					})

					// 取最晚的 10 个 messageBlock
					recentBlocks := messageBlock[len(messageBlock)-MaxBlocks:]

					// 将这些 messageBlock 的内容追加到 histories
					for _, mb := range recentBlocks {
						histories = append(histories, mb.Message...)
					}
				}
			}

			// lastMessage 是用户最后发送的消息，将它添加到 history 末尾
			// 为什么要这么做呢，因为最后个消息是不能单独成一个 block 的
			// Update: 这里应该做个判断，如果 lastMessage 和当前的 Message 一样，则不用添加
			if lastChatMessage.Id != histories[len(histories)-1].Id {
				histories = append(histories, lastChatMessage)
			}

			// 这样就能取到了剪裁后的，但是这样换汤不换药，之后还是得在 service 层面做分页。
			// 更新：已经做好啦
			err = u.llmService.StreamChat(c, llmChat, histories)
			if err != nil {
				u.logger.Sugar.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}()
	} else {
		response.Status(http.StatusInternalServerError).Error(consts.ErrChatStreamNotOpen).Send()
		return
	}

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
				Role:    schema.RoleToolCall,
				Content: "",
				ChatId:  chatEntity.Id,
			}

			if msg.Internal != nil && msg.Internal.ToolCall != nil {
				cm.ToolCall = (*schema.ToolCall)(msg.Internal.ToolCall)
			}

			messageList = append(messageList, cm)

			return true
		case schema.StateToolResponse:
			var cm = entity.ChatMessage{
				Role:    schema.RoleTool,
				Content: msg.ToolResponseMessage.Content,
				ChatId:  chatEntity.Id,
			}

			if msg.Internal != nil && msg.Internal.ToolCall != nil {
				cm.ToolCall = (*schema.ToolCall)(msg.Internal.ToolCall)
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
			return false
		case schema.StateFinished:
			tokenUsage = msg.TokenUsage
			return false
		case schema.StateToolFailed:
			// 这样会发生悬垂
			// 所以要添加个新的消息
			var cm = entity.ChatMessage{
				Role:    schema.RoleTool,
				Content: msg.ToolResponseMessage.Content,
				ChatId:  chatEntity.Id,
			}

			if msg.Internal != nil && msg.Internal.ToolCall != nil {
				cm.ToolCall = (*schema.ToolCall)(msg.Internal.ToolCall)
			}

			messageList = append(messageList, cm)

			return false
		default:
			return true
		}
	})

	// 发送 [DONE]
	c.SSEvent(eventName, eventDone)

	// close sse stream
	c.Writer.Flush()

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

	defer func() {
		// 移除缓存
		u.redis.Client.Del(c, streamIdCacheKey)
		u.redis.Client.Del(c, u.getCacheKey("entity:"+csc.ChatId.String()))
		u.redis.Client.Del(c, u.getCacheKey("stream:"+chatStreamRequest.StreamId+":user"))
		u.redis.Client.Del(c, chatIdStreamKey)
	}()
}

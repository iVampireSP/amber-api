package v1

import (
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"rag-new/internal/entity"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"strconv"
	"time"
)

// OpenAIChatCompletion godoc
// @Summary      OpenAI Chat Completion
// @Description  兼容 OpenAI Chat Completion 接口，认证需要使用 Assistant Share Token
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        chat  body  	schema.OpenAIChatCompletionRequest  true  "Chat"
// @Success      200  {object}  schema.ResponseBody{data=schema.OpenAIChatCompletionResponse}
// @Failure      400  {object}  schema.ResponseBody
// @Failure      500  {object}  schema.ResponseBody
// @Router       /api/openai-compatible/v1/chat/completions [post]
func (u *ChatController) OpenAIChatCompletion(c *gin.Context) {
	var response = schema.NewResponse(c)
	assistantEntity := u.assistantService.GetAssistantFromCtx(c)

	chatRequest := schema.OpenAIChatCompletionRequest{}
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		response.Status(http.StatusBadRequest).Error(err).Send()
		return
	}

	tools, err := u.assistantService.ToLLMTool(c, assistantEntity)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var prompt = u.getPrompt(assistantEntity, nil)
	var llmResponseChan = make(chan *schema.AssistantResponse)

	var llmChat = &schema.LLMChat{
		ResponseChan:   llmResponseChan,
		SystemPrompt:   prompt,
		UserPublicInfo: nil,
		Tools:          tools,
	}

	var histories = make([]*entity.ChatMessage, 0)
	// 转换
	for _, message := range chatRequest.Messages {
		histories = append(histories, &entity.ChatMessage{
			Role:    schema.ChatRole(message.Role),
			Content: message.Content,
		})
	}

	go func() {
		err = u.llmService.StreamChat(llmChat, histories)
		if err != nil {
			u.logger.Sugar.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()

	var created = time.Now().Unix()

	var tokenUsage = &schema.TokenUsage{
		PromptTokens:     0,
		CompletionTokens: 0,
		TotalTokens:      0,
	}

	var fakeChatId = "chatcmpl-" + strconv.FormatInt(created, 10)

	if chatRequest.Stream {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		c.Stream(func(w io.Writer) bool {
			// Emit Server Sent Events compatible
			msg, ok := <-llmResponseChan
			if !ok {
				return false
			}

			if msg == nil {
				return true
			}

			var openAIChatCompletionStreamResponse = schema.OpenAIChatCompletionStreamResponse{
				ID:      fakeChatId,
				Object:  "chat.completion.chunk",
				Created: created,
				Choices: []schema.OpenAIChatCompletionStreamResponseChoice{},
				Usage:   nil,
			}

			switch msg.State {
			case schema.StateChunk:
				openAIChatCompletionStreamResponse.Choices = []schema.OpenAIChatCompletionStreamResponseChoice{
					{
						Delta: schema.OpenAIChatCompletionStreamDelta{
							Content: msg.Content,
							Role:    string(schema.RoleAssistant),
						},
						Index:        0,
						FinishReason: nil,
					},
				}
				j, err := sonic.MarshalString(openAIChatCompletionStreamResponse)
				if err != nil {
					u.logger.Sugar.Error(err)
				}
				c.SSEvent(eventName, j)

				c.Writer.Flush()

				return true
			case schema.StateDone:
				tokenUsage = msg.TokenUsage
				return true
			case schema.StateFailed:
				return false
			case schema.StateFinished:
				return false
			case schema.StateToolFailed:
				return false
			default:
				return true
			}
		})

		j, err := sonic.MarshalString(schema.OpenAIChatCompletionStreamResponse{
			ID:      fakeChatId,
			Object:  "chat.completion.chunk",
			Created: created,
			Choices: []schema.OpenAIChatCompletionStreamResponseChoice{
				{
					Delta: schema.OpenAIChatCompletionStreamDelta{
						Content: "",
						Role:    string(schema.RoleAssistant),
					},
					Index:        0,
					FinishReason: "stop",
				},
			},
			Usage: nil,
		})
		if err != nil {
			u.logger.Sugar.Error(err)
		}
		c.SSEvent(eventName, j)

		j, err = sonic.MarshalString(schema.OpenAIChatCompletionStreamResponse{
			ID:      fakeChatId,
			Object:  "chat.completion.chunk",
			Created: created,
			Choices: []schema.OpenAIChatCompletionStreamResponseChoice{},
			Usage:   tokenUsage,
			Model:   chatRequest.Model,
		})
		if err != nil {
			u.logger.Sugar.Error(err)
		}
		c.SSEvent(eventName, j)

		c.SSEvent(eventName, eventDone)
		c.Writer.Flush()

		return
	}

	// 非 stream 模式
	var llmFullResponse = ""
	c.Stream(func(w io.Writer) bool {
		msg, ok := <-llmResponseChan
		if !ok {
			return false
		}
		if msg == nil {
			return true
		}
		switch msg.State {
		case schema.StateChunk:
			llmFullResponse += msg.Content
			return true
		case schema.StateDone:
			tokenUsage = msg.TokenUsage
			return true
		case schema.StateFailed:
			return false
		case schema.StateFinished:
			return false
		case schema.StateToolFailed:
			return false
		default:
			return true
		}
	})

	response.Data(schema.OpenAIChatCompletionResponse{
		ID:      fakeChatId,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   chatRequest.Model,
		Choices: []schema.OpenAIChatCompletionResponseChoice{
			{
				Message: schema.OpenAIChatCompletionRequestMessage{
					Content: llmFullResponse,
					Role:    string(schema.RoleAssistant),
				},
				Index: 0,
			},
		},
		Usage: tokenUsage,
	}).WithoutWrap().Error(err).Send()
}

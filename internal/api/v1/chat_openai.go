package v1

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"rag-new/internal/entity"
	_ "rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strconv"
	"strings"
	"time"
)

const maxImageSize = 3 * 1024 * 1024 // 3MB

// OpenAIChatCompletion godoc
// @Summary      OpenAI Chat Completion
// @Description  兼容 OpenAI Chat Completion 接口，认证需要使用 Assistant Share Token
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param 	     X-User-IP  header  string  false  "指定聊天中的用户 IP 地址，不指定则自动获取。此 IP 地址只会增加至 Prompt 中，如果不希望增加，请关闭系统自带 Prompt 选项"
// @Param        chat  body  	schema.OpenAIChatCompletionRequest  true  "Chat"
// @Success      200  {object}  schema.OpenAIChatCompletionResponse
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

	prompt, err := u.getPrompt(c, assistantEntity, nil, schema.OwnerGuest)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Send()
		return
	}

	var llmResponseChan = make(chan *schema.AssistantResponse)

	var llmChat = &schema.LLMChat{
		ResponseChan:   llmResponseChan,
		SystemPrompt:   prompt,
		UserPublicInfo: nil,
		Tools:          tools,
		Model:          chatRequest.Model,
	}

	// 如果不能使用，则用 auto
	if !u.config.OpenAI.CanUse(chatRequest.Model) {
		chatRequest.Model = consts.AutoModel
	}

	var histories = make([]*entity.ChatMessage, 0)
	// 转换
	for _, message := range chatRequest.Messages {
		var role = schema.ChatRole(message.Role)
		var content string
		var imageUrl string

		// 判断 Content 字段的类型
		switch contentTyped := message.Content.(type) {
		case string:
			content = contentTyped
		case []interface{}:
			// 循环每一个
			if len(contentTyped) > 0 {
				for _, v := range contentTyped {
					if textTyped, ok := v.(map[string]interface{}); ok {
						if textType, ok := textTyped["type"].(string); ok {
							// 如果是 text
							if textType == "text" {
								content = textTyped["text"].(string)
								if content == "" {
									response.Error(consts.ErrTextCannotBeEmpty).Send()
									return
								}
							} else if textType == "image_url" {
								// 读取下面的 image_url
								if imageUrlTyped, ok := textTyped["image_url"].(map[string]interface{}); ok {
									if imageUrlValue, ok := imageUrlTyped["url"].(string); ok {
										if imageUrlValue == "" {
											response.Error(consts.ErrImageUrlCannotBeEmpty).Send()
											return
										}

										imageUrl = imageUrlValue
									} else {
										response.Error(consts.ErrImageIsRequired).Send()
									}
								}
							}
						} else {
							response.Status(http.StatusBadRequest).Error(consts.ErrTypeRequired).Send()
							return
						}
					}
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content type"})
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content type"})
		}

		if imageUrl != "" {
			var file = &entity.File{}

			// 先检测否是 url(http)
			// 检测是否是 http 开头
			if strings.HasPrefix(imageUrl, "http") {
				file, err = u.fileService.CreateFileFromUrl(c, imageUrl)
				if err != nil {
					response.Status(http.StatusInternalServerError).Error(consts.ErrFileUrlNotURL).Send()
					return
				}
			} else {
				var substr = ";base64,"
				// 寻找 base64,，如果找到的话，去除前面的以及 base64,
				base64Index := strings.Index(imageUrl, ";base64,")
				if base64Index != -1 {
					imageUrl = imageUrl[base64Index+len(substr):]
				}

				// 如果 imageUrl 是 base64
				reader, err := base64ToReader(imageUrl)
				if err != nil {
					response.Status(http.StatusInternalServerError).Error(consts.ErrFileUrlNotValidBase64).Send()
					return
				}

				// 如果是 nil
				if reader == nil {
					response.Status(http.StatusInternalServerError).Error(consts.ErrFileUrlNotValidBase64).Send()
					return
				}

				file, err = u.fileService.CreateFile(c, reader)
				if err != nil {
					response.Status(http.StatusInternalServerError).Error(err).Send()
					return
				}
			}

			// 如果是文件的话，要新增两次
			histories = append(histories, &entity.ChatMessage{
				Role:    schema.RoleFile,
				Content: file.Id.String(),
			})
			histories = append(histories, &entity.ChatMessage{
				Role:    role,
				Content: content,
			})

		} else {
			histories = append(histories, &entity.ChatMessage{
				Role:    role,
				Content: content,
			})
		}

	}

	go func() {
		err = u.llmService.StreamChat(c, llmChat, histories)
		if err != nil {
			u.logger.Sugar.Error(err)
			response.Status(http.StatusBadRequest).Error(err).Send()
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
				Model:   chatRequest.Model,
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

				_, err = c.Writer.WriteString(eventName + ": " + j + "\n\n")
				if err != nil {
					u.logger.Sugar.Error(err)
					return false
				}

				c.Writer.Flush()

				return true
			case schema.StateDone:
				return true
			case schema.StateFailed:
				tokenUsage = msg.TokenUsage
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
			Model: chatRequest.Model,
			Usage: nil,
		})
		if err != nil {
			u.logger.Sugar.Error(err)
		}
		// 直接 write
		_, err = c.Writer.WriteString(eventName + ": " + j + "\n\n")
		if err != nil {
			u.logger.Sugar.Error(err)
			return
		}

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
		_, err = c.Writer.WriteString(eventName + ": " + j + "\n\n")
		if err != nil {
			u.logger.Sugar.Error(err)
			return
		}

		_, err = c.Writer.WriteString(eventName + ": " + eventDone + "\n\n")
		if err != nil {
			u.logger.Sugar.Error(err)
			return
		}
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

func base64ToReader(s string) (io.ReadSeeker, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	if len(decoded) > maxImageSize {
		return nil, fmt.Errorf("image size exceeds maximum allowed size of %d bytes", maxImageSize)
	}
	return bytes.NewReader(decoded), nil
}

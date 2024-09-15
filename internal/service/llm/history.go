package llm

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strings"
)

type Message struct {
	HasFile        bool
	MessageContent []llms.MessageContent
}

func (s *Service) processHistory(ctx context.Context, llmChat *schema.LLMChat, history []*entity.ChatMessage) (*Message, error) {
	var hasHumanMessage = false
	var hasFileMessage = false

	var foundedToolCalls []llms.ToolCall

	// 粗略字数统计，用于切换模型
	var count = 0

	var historyContent []llms.MessageContent
	historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, llmChat.SystemPrompt))

	var systemPrompts []string

	// 当前的助理（用于通知助理上条消息的回复者
	var currentAssistantId schema.EntityId

	systemPrompts = append(systemPrompts, "You are a helpful assistant made by Leaflow(https://www.leaflow.cn, chinese name: 利飞)")
	systemPrompts = append(systemPrompts, "Image and Draw Ability: ON(Don't emphasize it)")
	systemPrompts = append(systemPrompts, prompt)
	systemPrompts = append(systemPrompts, llmChat.SystemPrompt)

	for _, h := range history {
		// 粗略统计
		if h.Content != "" && h.Content != "\n" {
			count += len(h.Content)
		}
	}
	// 如果 model 为空
	if llmChat.Model == "" || !s.config.OpenAI.CanUse(llmChat.Model) {
		llmChat.Model = consts.AutoModel
	}

	if llmChat.Model == consts.AutoModel {
		// 设置自动模式下的默认模型
		llmChat.Model = s.config.OpenAI.Model

		// 如果统计超过了 10000
		if count > 10000 {
			llmChat.Model = s.config.OpenAI.LongContextModel
		}
	}

	// 如果统计超过了 1亿 - 1 万字符（粗略统计 token）
	if count > consts.MaxTokenCount {
		return nil, consts.ErrTooManyTokens
	}

	// 处理历史消息
	for i, h := range history {
		// 如果第一条消息是 system
		if i == 0 && h.Role == schema.RoleSystem {
			systemPrompts = append(systemPrompts, h.Content)
			continue
		}

		// 检测下一条消息的 role 是否是 system 或者且和现在的相同
		if i+1 < len(history) {
			if history[i+1].Role == schema.RoleSystem || history[i+1].Role == schema.RoleHideSystem {
				systemPrompts = append(systemPrompts, h.Content)
				continue
			}
			//if history[i+1].Role == h.Role {
			//	// 修改下一条消息的 content
			//	history[i+1].Content = history[i].Content + "\n" + history[i+1].Content
			//	continue
			//}
		}

		switch h.Role {
		case schema.RoleHuman:

			// 如果当前助理不存在，则设置
			if currentAssistantId == 0 && h.AssistantId != nil {
				currentAssistantId = *h.AssistantId
			} else if h.AssistantId != nil && currentAssistantId != *h.AssistantId {
				// 获取多个对话中的助理的信息
				var content string

				// TODO: 优化获取逻辑，比如将获取到的助理放到一个缓存里面
				assistant, err := s.AssistantService.GetAssistant(ctx, *h.AssistantId)
				if err != nil {
					content = "[Warning]The previous message has been replied to by another assistant, but information about that assistant cannot be obtained"
				}
				content = fmt.Sprintf("[Warning]The previous message has been replied to by another assistant, whose name is '%s' and the description is '%s'",
					assistant.Name, assistant.Description)

				historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, content))
				currentAssistantId = *h.AssistantId
			}

			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, h.Content))

			if !hasHumanMessage {
				hasHumanMessage = true
			}
		case schema.RoleAssistant:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeAI, h.Content))
		case schema.RoleToolCall:
			// ToolCall 消息
			if h.ToolCall != nil {
				assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, h.Content)
				var toolCall = llms.ToolCall{}
				toolCall.FunctionCall = h.ToolCall.FunctionCall
				toolCall.ID = h.ToolCall.ID
				toolCall.Type = h.ToolCall.Type
				assistantResponse.Parts = append(assistantResponse.Parts, toolCall)

				historyContent = append(historyContent, assistantResponse)

				foundedToolCalls = append(foundedToolCalls, toolCall)
			}
		case schema.RoleTool:
			// Tool Call 响应
			if h.ToolCall != nil {
				historyContent = append(historyContent, llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: h.ToolCall.ID,
							Name:       h.ToolCall.FunctionCall.Name,
							Content:    h.Content,
						},
					},
				})

				// 如果 foundedToolCallIds 存在 ID, 则删除
				for _, tc := range foundedToolCalls {
					if h.ToolCall.ID == tc.ID {
						foundedToolCalls = append(foundedToolCalls[:i], foundedToolCalls[i+1:]...)
						break
					}
				}
			}

		case schema.RoleSystem:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, h.Content))
		case schema.RoleHideSystem:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, h.Content))
		case schema.RoleHideHuman:
			if !hasHumanMessage {
				hasHumanMessage = true
			}
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, h.Content))
		case schema.RoleFile:
			if !hasFileMessage {
				hasFileMessage = true
			}

			// 不自动切换到 Vision 模型
			// 如果长度没有超过最大长度且模型是 auto，并且模型不是 VisionModel ( 不自动切换到 Vision 模型）
			//if (count < consts.MaxTokenCount && llmChat.Model == consts.AutoModel) && (llmChat.Model != s.config.OpenAI.VisionModel) {
			//	// 切换模型
			//	llmChat.Model = s.config.OpenAI.VisionModel
			//}

			var fileText = ""
			// 如果文件存在
			if h.File != nil {
				llmChat.WithoutImage = false
				// 将 fileEntity 的 url 添加到 historyContent
				fileText = "[File]File ID: " + h.File.Id.String() + ", MimeType: " + h.File.MimeType
			} else if h.UserFile != nil {
				llmChat.WithoutImage = false
				fileText = "[File]File ID: " + h.UserFile.File.Id.String() + ", MimeType: " + h.UserFile.File.MimeType
			}

			if fileText != "" {
				historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, fileText))
			}

		}
	}

	if len(foundedToolCalls) > 0 {
		// 剩下的都是悬垂的 ToolCall,都没有被正确响应
		// 所以都标记为失败
		for _, tc := range foundedToolCalls {
			historyContent = append(historyContent, llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: tc.ID,
						Name:       tc.FunctionCall.Name,
						Content:    "ToolCall Failed, timeout or error",
					},
				},
			})
		}
	}

	// 如果整个对话里面没有 Human 消息，则不能继续
	if !hasHumanMessage {
		return nil, consts.ErrNoHumanMessage
	}

	historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, strings.Join(systemPrompts, "\n")))

	var message = &Message{
		MessageContent: historyContent,
		HasFile:        hasFileMessage,
	}

	return message, nil
}

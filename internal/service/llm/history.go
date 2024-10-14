package llm

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/builtin_tool"
	"rag-new/pkg/consts"
	"strings"
)

type Message struct {
	HasFile        bool
	MessageContent []llms.MessageContent
}

const defaultToolFailed = "ToolCall Failed, timeout or error"

func (s *Service) processHistory(_ context.Context, llmChat *schema.LLMChat, history []*entity.ChatMessage) (*Message, error) {
	var hasHumanMessage = false
	var hasFileMessage = false

	var lastToolCall *llms.ToolCall

	var historyContent []llms.MessageContent

	var systemPrompts []string

	// 如果没有禁用 Agent 方式的图片工具
	if !llmChat.WithoutImage {
		systemPrompts = append(systemPrompts, "Image and Draw Ability: ON(Don't emphasize it)")
	}

	if !llmChat.WithoutBrowsing {
		systemPrompts = append(systemPrompts, builtin_tool.WebSearchPrompt)
	}

	// 当前的助理（用于通知助理上条消息的回复者
	var currentAssistantId schema.EntityId

	systemPrompts = append(systemPrompts, prompt)
	systemPrompts = append(systemPrompts, llmChat.SystemPrompt)

	var lastKnowledgeMessage string

	// 如果 model 为空
	if llmChat.Model == "" || !s.config.OpenAI.CanUse(llmChat.Model) {
		llmChat.Model = consts.AutoModel
	}

	if llmChat.Model == consts.AutoModel {

		// 粗略字数统计，用于切换模型
		var count = 0
		for _, h := range history {
			// 粗略统计
			if h.Content != "" && h.Content != "\n" {
				count += len(h.Content)
			}
		}

		// 设置自动模式下的默认模型
		llmChat.Model = s.config.OpenAI.Model

		// 如果统计超过了 10000
		if count > 100000 {
			llmChat.Model = s.config.OpenAI.LongContextModel
		}

		// 如果统计超过了 1亿 - 1 万字符（粗略统计 token）
		if count > consts.MaxTokenCount {
			return nil, consts.ErrTooManyTokens
		}
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

		var timeString = ""
		// 创建时间，如果 h.CreatedAt 已设置（不为 0000）
		if !h.CreatedAt.IsZero() {
			// 将创建时间转换为字符串
			timeString = fmt.Sprintf("%s", h.CreatedAt.Format("2006-01-02 15:04:05"))
		}

		switch h.Role {
		// RoleSystem 和 RoleHideSystem 为系统消息
		case schema.RoleSystem:
			if h.Content == "" {
				continue
			}

			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, h.Content))
		case schema.RoleHideSystem:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, h.Content))
		//	RoleHuman，RoleHideHuman，RoleHumanLater 是用户消息
		case schema.RoleHuman:

			// 获取多个对话中的助理的信息
			// 如果当前助理不存在，则设置
			if currentAssistantId == 0 && h.AssistantId != nil {
				currentAssistantId = *h.AssistantId
			}

			if timeString != "" {
				h.Content = fmt.Sprintf("[Sent at %s]\n%s", timeString, h.Content)
			}

			newContent, err := s.optimizeFlow(h.Content)
			if err != nil {
				return nil, err
			}
			h.Content = newContent

			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, h.Content))

			if !hasHumanMessage {
				hasHumanMessage = true
			}
		case schema.RoleHideHuman:
			if !hasHumanMessage {
				hasHumanMessage = true
			}

			newContent, err := s.optimizeFlow(h.Content)
			if err != nil {
				return nil, err
			}
			h.Content = newContent

			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, h.Content))

		case schema.RoleHumanLater:
			newContent, err := s.optimizeFlow(h.Content)
			if err != nil {
				return nil, err
			}
			h.Content = newContent

			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, h.Content))

		case schema.RoleAssistant:
			// 检测是否是悬垂调用
			if lastToolCall != nil {
				// 上条消息可能有问题，将上个 ToolCall 标记为失败
				historyContent = append(historyContent, llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: lastToolCall.ID,
							Name:       lastToolCall.FunctionCall.Name,
							Content:    defaultToolFailed,
						},
					},
				})
				lastToolCall = nil
			}

			var systemContent = ""
			// 也不一定不存在，因为可能上个消息没有助理
			//if h.AssistantId == nil {
			//	// 这说明上个助理不存在
			//	systemContent = "[Warning]The previous message has been replied to by another assistant, no you, but can not get assistant info."
			//}
			if h.Assistant != nil {
				systemContent = fmt.Sprintf("[Warning]The previous message has been replied to by another assistant, whose name is '%s' and the description is '%s', replid at %s",
					h.Assistant.Name, h.Assistant.Description, timeString)

				currentAssistantId = *h.AssistantId

			}

			if systemContent != "" {
				historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, systemContent))
			}

			if timeString != "" {
				h.Content = fmt.Sprintf("[Sent at %s]\n%s", timeString, h.Content)
			}

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

				// 因为 ToolCall 消息的下一条消息必须是 tool
				lastToolCall = &toolCall
			}
		case schema.RoleTool:
			// Tool Call 响应
			if h.ToolCall != nil {
				if lastToolCall != nil && lastToolCall.ID == h.ToolCall.ID {
					lastToolCall = nil

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
				} else if lastToolCall != nil {
					// 不相同，说明有问题,将上个 Tool Call 标记为失败
					historyContent = append(historyContent, llms.MessageContent{
						Role: llms.ChatMessageTypeTool,
						Parts: []llms.ContentPart{
							llms.ToolCallResponse{
								ToolCallID: lastToolCall.ID,
								Name:       lastToolCall.FunctionCall.Name,
								Content:    defaultToolFailed,
							},
						},
					})
					lastToolCall = nil
				}

			}

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

			if h.File != nil {
				// 如果是 Vision 模型，直接处理
				if llmChat.UseVisionModel {
					// 禁用图片解释工具
					llmChat.WithoutImage = true

					// 直接添加图片
					fileUrl, err := s.FileService.GetImageUrl(h.File)
					if err != nil {
						return nil, err
					}

					// 获取下一条消息，如果 i+1 有内容且 role 为 Human
					if i+1 < len(history) && history[i+1].Role == schema.RoleHuman {
						// 获取下一条消息
						nextMessage := history[i+1]
						historyContent = append(historyContent, llms.MessageContent{
							Role: llms.ChatMessageTypeHuman,
							Parts: []llms.ContentPart{
								llms.ImageURLWithDetailPart(fileUrl, "auto"),
								llms.TextPart(nextMessage.Content),
							},
						})
					}
				} else {
					// 如果不是，则只添加文件信息，交给 Agent 处理
					// 将 fileEntity 的 url 添加到 historyContent
					fileText := "[Upload File]File Hash: " + h.File.FileHash + ", MimeType: " + h.File.MimeType + ", "
					historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, fileText))
				}
			}

		case schema.RoleKnowledge:
			if h.Content == "" {
				h.Content = consts.LibraryResultEmptyPrompt
			} else {
				h.Content = consts.LibraryResultPrompt + "\n" + h.Content
			}

			lastKnowledgeMessage = h.Content
		}
	}

	// 将知识库消息放到 system 消息里面
	if lastKnowledgeMessage != "" {
		systemPrompts = append(systemPrompts, lastKnowledgeMessage)
	}

	// 如果整个对话里面没有 Human 消息，则不能继续
	if !hasHumanMessage {
		return nil, consts.ErrNoHumanMessage
	}

	// 拼接系统 Prompt 并放入最底
	historyContent = append(historyContent,
		llms.TextParts(llms.ChatMessageTypeSystem, strings.Join(systemPrompts, "\n")))

	var message = &Message{
		MessageContent: historyContent,
		HasFile:        hasFileMessage,
	}

	return message, nil
}

type flowFunc func(content string) (string, error)

func (s *Service) optimizeFlow(content string) (string, error) {
	var flow = []flowFunc{
		s.flowCleanupNewLine,
	}

	for _, f := range flow {
		var err error
		content, err = f(content)
		if err != nil {
			return "", err
		}
	}

	return content, nil
}

func (s *Service) flowCleanupNewLine(content string) (string, error) {
	// 匹配连续的 \r\n
	mergedText := strings.ReplaceAll(content, "\n\n", "\n")

	// 继续合并，直到没有连续的换行
	for strings.Contains(mergedText, "\n\n") {
		mergedText = strings.ReplaceAll(mergedText, "\n\n", "\n")
	}

	return mergedText, nil
}

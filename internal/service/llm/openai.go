package llm

import (
	"bytes"
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/tmc/langchaingo/llms"
	"io"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strconv"
	"strings"
)

// StreamChat 执行对话
func (s *Service) StreamChat(llmChat *schema.LLMChat, history []*entity.ChatMessage) error {
	var historyContent []llms.MessageContent

	var hasHumanMessage = false

	historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, llmChat.SystemPrompt))

	for i, h := range history {
		// 如果第一条消息是 system
		if i == 0 && h.Role == schema.RoleSystem {
			var newPrompt = ""
			if llmChat.SystemPrompt == "" {
				newPrompt = h.Content
			} else {
				newPrompt = llmChat.SystemPrompt + "\n" + h.Content
			}
			historyContent[0] = llms.TextParts(llms.ChatMessageTypeSystem, newPrompt)
			continue
		}

		// 检测下一条消息的 role 是否是 system 或者且和现在的相同
		if i+1 < len(history) {
			if history[i+1].Role == schema.RoleSystem || history[i+1].Role == schema.RoleHideSystem {
				history[i+1].Content = history[i].Content + "\n" + history[i+1].Content
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
			//content := "[User says] " + h.Content
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, h.Content))

			if !hasHumanMessage {
				hasHumanMessage = true
			}
		case schema.RoleAssistant:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeAI, h.Content))
		case schema.RoleSystem:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, h.Content))
		case schema.RoleHideSystem:
			//content := "[System Hint]" + h.Content
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, h.Content))
		}
	}

	if !hasHumanMessage {
		return consts.ErrNoHumanMessage
	}

	var requestAgain = true

	for {
		var fullResponse [][]byte

		ctx := context.Background()

		//fmt.Println("再次请求吗?" + fmt.Sprint(requestAgain))
		if !requestAgain {
			break
		}

		requestAgain = false
		var isToolCall = false

		//fmt.Println("对话历史", historyContent)

		resp, err := s.OpenAI.GenerateContent(ctx,
			historyContent,
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				// 检测长度
				if len(chunk) == 0 {
					return nil
				}

				//fmt.Printf("Received chunk: %s\n", chunk)

				// 检测是否 json
				var isJson = sonic.Valid(chunk)
				if !isJson {
					var stringChunk = string(chunk)

					llmChat.ResponseChan <- &schema.AssistantResponse{
						State: schema.StateChunk,
						ChunkMessage: &schema.ChunkMessage{
							Content: stringChunk,
						},
						Content: stringChunk,
					}

					fullResponse = append(fullResponse, chunk)
				}

				return nil
			}),
			llms.WithTools(llmChat.Tools),
			llms.WithN(llmChat.N),
			llms.WithMaxTokens(llmChat.MaxTokens),
			llms.WithTemperature(llmChat.Temperature),
			llms.WithTopP(llmChat.TopP),
			llms.WithTopK(llmChat.TopK))
		if err != nil {
			llmChat.ResponseChan <- &schema.AssistantResponse{
				State:   schema.StateFailed,
				Content: err.Error(),
			}
			return err
		}

		respChoice := resp.Choices[0]
		tokenUsage := s.getTokenUsage(respChoice)

		if respChoice.FuncCall != nil {
			//fmt.Println("FunCall 检测到工具调用")

			isToolCall = true
			requestAgain = true
		}

		if isToolCall {
			//fmt.Printf("正在调用: %v\n", respChoice.FuncCall.Name)

			var fullArgs = ""

			assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, respChoice.Content)
			for _, tc := range respChoice.ToolCalls {
				// 拼接参数
				fullArgs += tc.FunctionCall.Arguments
			}

			var firstToolCall = respChoice.ToolCalls[0]
			firstToolCall.FunctionCall.Arguments = fullArgs
			assistantResponse.Parts = append(assistantResponse.Parts, firstToolCall)

			historyContent = append(historyContent, assistantResponse)

			// 去除 fullArgs 的首尾 \n（一直检测）
			for {
				if fullArgs[0] == '\n' {
					fullArgs = fullArgs[1:]
				} else if fullArgs[len(fullArgs)-1] == '\n' {
					fullArgs = fullArgs[:len(fullArgs)-1]
				} else {
					break
				}
			}

			// 解析工具
			var functionCallArgs schema.FunctionCallArgs
			err = sonic.Unmarshal([]byte(fullArgs), &functionCallArgs)
			if err != nil {
				return err
			}

			tool, functionName, err := s.spiltFunctionName(respChoice.FuncCall.Name)
			if err != nil {
				llmChat.ResponseChan <- &schema.AssistantResponse{
					State:      schema.StateFailed,
					Content:    err.Error(),
					TokenUsage: tokenUsage,
				}
				return err
			}

			llmChat.ResponseChan <- &schema.AssistantResponse{
				State: schema.StateToolCalling,
				ToolCallMessage: &schema.ToolCallMessage{
					ToolName:     tool.Name,
					FunctionName: respChoice.FuncCall.Name,
					Args:         functionCallArgs,
				},
				TokenUsage: tokenUsage,
			}

			remoteFunctionResponse, err := s.callRemoteFunction(tool, llmChat.UserPublicInfo, functionName, functionCallArgs)
			if err != nil {
				llmChat.ResponseChan <- &schema.AssistantResponse{
					State:   schema.StateToolFailed,
					Content: err.Error(),
					ToolResponseMessage: &schema.ToolResponseMessage{
						ToolName:     tool.Name,
						FunctionName: respChoice.FuncCall.Name,
						Content:      err.Error(),
					},
					TokenUsage: tokenUsage,
				}
				return err
				//remoteFunctionResponse.Content = err.Error()
			} else {
				historyContent = append(historyContent, llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: respChoice.ToolCalls[0].ID,
							Name:       respChoice.FuncCall.Name,
							Content:    remoteFunctionResponse.Content,
						},
					},
				})

				llmChat.ResponseChan <- &schema.AssistantResponse{
					State: schema.StateToolResponse,
					ToolResponseMessage: &schema.ToolResponseMessage{
						ToolName:     tool.Name,
						FunctionName: respChoice.FuncCall.Name,
						Content:      remoteFunctionResponse.Content,
					},
					TokenUsage: tokenUsage,
				}
			}

		} else {
			requestAgain = false
		}

		llmChat.ResponseChan <- &schema.AssistantResponse{
			State:      schema.StateDone,
			TokenUsage: tokenUsage,
		}

		historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))

		//fmt.Println("本轮历史", historyContent)
	}

	llmChat.ResponseChan <- &schema.AssistantResponse{
		State: schema.StateFinished,
	}

	return nil
}

func (s *Service) spiltFunctionName(functionName string) (*entity.Tool, string, error) {
	// 根据 _ 分割
	var functionNames = strings.Split(functionName, "_")

	// 从第 1 个开始取到最后一个
	var toolName = strings.Join(functionNames[1:], "_")

	// 第一个是 id
	toolId, err := strconv.Atoi(functionNames[0])
	if err != nil {
		return nil, toolName, err
	}

	// 从数据库中获取
	tool, err := s.ToolService.GetTool(context.Background(), int64(toolId))

	return tool, toolName, err
}

func (s *Service) callRemoteFunction(tool *entity.Tool, userPublicInfo *schema.UserPublicInfo, functionName string, args schema.FunctionCallArgs) (*schema.ToolRemoteResponse, error) {
	var callbackUrl = tool.Data.CallbackUrl

	var toolRequest = &schema.ToolRemoteRequest{
		FunctionName: functionName,
		Parameters:   args,
		ApiKey:       tool.ApiKey,
	}

	if userPublicInfo != nil {
		toolRequest.User = userPublicInfo
	}

	toolRequestJson, err := sonic.Marshal(toolRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", callbackUrl, bytes.NewBuffer(toolRequestJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+toolRequest.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyJson := &schema.ToolRemoteResponse{}

	err = sonic.Unmarshal(body, bodyJson)
	if err != nil {
		return nil, err
	}

	if bodyJson.Success {
		return bodyJson, nil
	}

	return bodyJson, errors.New(bodyJson.Content)
}

func (s *Service) getTokenUsage(respChoice *llms.ContentChoice) *schema.TokenUsage {
	var tokenUsage = &schema.TokenUsage{}

	// 如果 respChoice.GenerationInfo 中有 prompt_tokens
	if respChoice.GenerationInfo["PromptTokens"] != nil {
		tokenUsage.PromptTokens = respChoice.GenerationInfo["PromptTokens"].(int)
	} else {
		tokenUsage.PromptTokens = 0
	}

	// 如果 respChoice.GenerationInfo 中有 completion_tokens
	if respChoice.GenerationInfo["CompletionTokens"] != nil {
		tokenUsage.CompletionTokens = respChoice.GenerationInfo["CompletionTokens"].(int)
	} else {
		tokenUsage.CompletionTokens = 0
	}

	// 如果 respChoice.GenerationInfo 中有 total_tokens
	if respChoice.GenerationInfo["TotalTokens"] != nil {
		tokenUsage.TotalTokens = respChoice.GenerationInfo["TotalTokens"].(int)
	} else {
		tokenUsage.TotalTokens = 0
	}

	return tokenUsage
}

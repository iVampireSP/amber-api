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
	"rag-new/internal/service/builtin_tool"
	"rag-new/pkg/consts"
	"strconv"
	"strings"
)

// 强制停止（如果连续函数调用超过 4 次，则强制停止输出）
const forceStopCount = 4

// 警告次数（如果 LLM 连续调用超过 3 次，则警告输出）
const warningCount = 3

// 警告 LLM 调用太多次工具， 要求停止
const warningMessage = "[Warning]You are attempting to call the tool/function repeatedly, please use the tool/function properly and stop response. If you continue to call repeatedly, the chat will be forcibly terminated."
const forceStopSystemMessage = "[Force Stop]You have still repeatedly called the tool/function many times, and the chat has been forcibly terminated."

// StreamChat 执行对话
func (s *Service) StreamChat(llmChat *schema.LLMChat, history []*entity.ChatMessage) error {
	// 不要从接受侧关闭 channel
	defer close(llmChat.ResponseChan)
	historyContent, err := s.processHistory(llmChat, history)
	if err != nil {
		return err
	}

	// 是否再次请求
	var requestAgain = true
	// 连续的函数调用次数
	var functionCallCount = 0

	for {
		ctx := context.Background()

		//fmt.Println("再次请求吗?" + fmt.Sprint(requestAgain))
		if !requestAgain {
			break
		}

		// 标记不再请求，因为默认情况是不需要调用工具的。如果需要调用工具并等待工具回应，则需要设置为 true
		requestAgain = false

		// 计算工具调用次数
		if functionCallCount >= forceStopCount {
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, forceStopSystemMessage))
			continue
		} else if functionCallCount >= warningCount {
			// 添加 system 消息
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, warningMessage))
		}

		var llmTools = append(s.BuiltInTools.GetTools(), llmChat.Tools...)

		resp, err := s.GenerateContent(ctx, llmChat, llmTools, historyContent)
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
			// 检测到工具调用，标记为需要再次请求
			requestAgain = true

			// 将计数器加一次
			functionCallCount += 1

			// 拼接完整参数
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
			var functionCallArgs schema.FunctionCallArguments
			err = sonic.Unmarshal([]byte(fullArgs), &functionCallArgs)
			if err != nil {
				return err
			}

			prefix, functionName := s.spiltFunctionName(respChoice.FuncCall.Name)
			//if err != nil {
			//	llmChat.ResponseChan <- &schema.AssistantResponse{
			//		State:      schema.StateFailed,
			//		Content:    err.Error(),
			//		TokenUsage: tokenUsage,
			//	}
			//	return err
			//}

			var toolCalling = &schema.AssistantResponse{
				State: schema.StateToolCalling,
				ToolCallMessage: &schema.ToolCallMessage{
					FunctionName: respChoice.FuncCall.Name,
					Arguments:    functionCallArgs,
				},
				TokenUsage: tokenUsage,
			}

			var toolRemoteResponse = &schema.ToolRemoteResponse{}
			var toolName = ""
			if prefix == builtin_tool.NAME {
				// 是 builtin，则调用内置函数
				toolRemoteResponse, err = s.BuiltInTools.CallFunction(ctx, functionName, functionCallArgs)
				toolName = builtin_tool.NAME
				if err != nil {
					// 也许内置函数不应该报 ToolFailed,不如直接 failed
					llmChat.ResponseChan <- &schema.AssistantResponse{
						State:   schema.StateFailed,
						Content: err.Error(),
						//ToolResponseMessage: &schema.ToolResponseMessage{
						//	ToolName:     builtin_tool.NAME,
						//	FunctionName: respChoice.FuncCall.Name,
						//	Content:      err.Error(),
						//},
						TokenUsage: tokenUsage,
					}
					return err
				}

				//toolCalling.ToolCallMessage.ToolName = builtin_tool.NAME
			} else {
				// 转换 prefix
				toolId, err := strconv.Atoi(prefix)
				if err != nil {
					llmChat.ResponseChan <- &schema.AssistantResponse{
						// 这里改成 failed 会不会更好？
						State:   schema.StateToolFailed,
						Content: err.Error(),
						ToolResponseMessage: &schema.ToolResponseMessage{
							ToolName:     builtin_tool.NAME,
							FunctionName: respChoice.FuncCall.Name,
							Content:      err.Error(),
						},
						TokenUsage: tokenUsage,
					}
					return err
				}

				// 获取 Tool
				selectedTool, err := s.GetToolById(ctx, int64(toolId))
				if err != nil {
					llmChat.ResponseChan <- &schema.AssistantResponse{
						// 这里改成 failed 会不会更好？
						State:   schema.StateToolFailed,
						Content: err.Error(),
						ToolResponseMessage: &schema.ToolResponseMessage{
							ToolName:     toolName,
							FunctionName: respChoice.FuncCall.Name,
							Content:      err.Error(),
						},
						TokenUsage: tokenUsage,
					}
					return err
				}

				toolName = selectedTool.Name

				// 调用远程函数
				toolRemoteResponse, err = s.callRemoteFunction(selectedTool, llmChat, functionName, functionCallArgs)
				if err != nil {
					llmChat.ResponseChan <- &schema.AssistantResponse{
						State:   schema.StateToolFailed,
						Content: err.Error(),
						ToolResponseMessage: &schema.ToolResponseMessage{
							ToolName:     toolName,
							FunctionName: respChoice.FuncCall.Name,
							Content:      err.Error(),
						},
						TokenUsage: tokenUsage,
					}
					return err
				}
			}

			// 如果是 builtin ，则不告知
			if toolName != builtin_tool.NAME {
				llmChat.ResponseChan <- toolCalling
				llmChat.ResponseChan <- &schema.AssistantResponse{
					State: schema.StateToolResponse,
					ToolResponseMessage: &schema.ToolResponseMessage{
						ToolName:         toolName,
						FunctionName:     respChoice.FuncCall.Name,
						Content:          toolRemoteResponse.Content,
						RememberResponse: toolRemoteResponse.RememberResponse,
						StopGeneration:   toolRemoteResponse.StopGeneration,
					},
					TokenUsage: tokenUsage,
				}
			}

			historyContent = append(historyContent, llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: respChoice.ToolCalls[0].ID,
						Name:       respChoice.FuncCall.Name,
						Content:    toolRemoteResponse.Content,
					},
				},
			})

			// 如果函数要求停止生成
			if toolRemoteResponse.StopGeneration {
				requestAgain = false
			}

		} else {
			// 不是工具调用，不再进行新的一轮请求，然后清除计数器
			requestAgain = false
			functionCallCount = 0
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

// spiltFunctionName 将函数名分割为 entity_toolName （entity 为 *entity.Tool，toolName 为 string）的形式
func (s *Service) spiltFunctionName(functionName string) (prefix string, realFunctionName string) {
	// 根据 _ 分割
	var functionNames = strings.Split(functionName, "_")

	// 从第 1 个开始取到最后一个
	var toolName = strings.Join(functionNames[1:], "_")

	return functionNames[0], toolName
}

func (s *Service) GetToolById(ctx context.Context, id int64) (*entity.Tool, error) {
	return s.ToolService.GetTool(ctx, id)
}

// callRemoteFunction 可以调用远程函数
func (s *Service) callRemoteFunction(tool *entity.Tool, llmChat *schema.LLMChat, functionName string, args schema.FunctionCallArguments) (*schema.ToolRemoteResponse, error) {
	if !s.config.Debug.Enabled {
		internalAddress, err := s.ToolService.IsAllowed(tool.Data.CallbackUrl)
		if err != nil {
			return nil, err
		}
		if internalAddress {
			return nil, consts.ErrToolAddressIsInternal
		}
	}

	var toolRequest = &schema.ToolRemoteRequest{
		FunctionName: functionName,
		Parameters:   args,
		Chat:         llmChat.Chat,
	}

	if llmChat.UserPublicInfo != nil {
		toolRequest.User = llmChat.UserPublicInfo
	}

	toolRequestJson, err := sonic.Marshal(toolRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tool.Data.CallbackUrl, bytes.NewBuffer(toolRequestJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if tool.ApiKey != "" {
		req.Header.Set("Authorization", "Bearer "+tool.ApiKey)
	}

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

func (s *Service) processHistory(llmChat *schema.LLMChat, history []*entity.ChatMessage) ([]llms.MessageContent, error) {
	var hasHumanMessage = false
	var hasImageMessage = false

	var historyContent []llms.MessageContent
	historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, llmChat.SystemPrompt))
	var systemPrompts []string

	systemPrompts = append(systemPrompts, "Image Ability: ON")
	systemPrompts = append(systemPrompts, llmChat.SystemPrompt)

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
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, h.Content))
		case schema.RoleHideHuman:
			if !hasHumanMessage {
				hasHumanMessage = true
			}
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, h.Content))
		case schema.RoleImage:
			if !hasImageMessage {
				hasImageMessage = true
			}
			var imageText = "[Image]Image ID: " + h.Content
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, imageText))
		}
	}

	// 如果整个对话里面没有 Human 消息，则不能继续
	if !hasHumanMessage {
		return historyContent, consts.ErrNoHumanMessage
	}

	var imagePrompt = `
The chat does not have images. you can't use built-in image tools. image_id can only get from user's uploaded images.
`

	// 如果有图片消息
	if hasImageMessage {
		imagePrompt = `
The chat has images, you can use built-in image tools.
`
	}

	systemPrompts = append(systemPrompts, imagePrompt)

	historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, strings.Join(systemPrompts, "\n")))

	return historyContent, nil
}

func (s *Service) GenerateContent(ctx context.Context, llmChat *schema.LLMChat, llmTools []llms.Tool, historyContent []llms.MessageContent) (response *llms.ContentResponse, err error) {
	// 上一个字
	var lastWord = ""
	// 重复次数
	var lastWordRepeatCount = 0

	resp, err := s.OpenAI.GenerateContent(ctx,
		historyContent,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			// 检测长度
			if len(chunk) == 0 {
				return nil
			}

			// 取 chunk 中最后一个字
			var chunkLastWord = string(chunk[len(chunk)-1])
			// 检测是否是上一个字
			if lastWord == chunkLastWord {
				lastWordRepeatCount++
			} else {
				lastWordRepeatCount = 0
				lastWord = chunkLastWord
			}
			// 如果上一个字重复次数大于 10，就终止
			if lastWordRepeatCount >= 10 {
				return consts.ErrWordRepeatedDetected
			}

			//fmt.Printf("Received chunk: %s\n", chunk)

			// 检测是否 json，判断是否是工具调用
			var isJson = sonic.Valid(chunk)
			if !isJson {
				var stringChunk = string(chunk)

				err = s.sendResponse(llmChat.ResponseChan, schema.AssistantResponse{
					State: schema.StateChunk,
					ChunkMessage: &schema.ChunkMessage{
						Content: stringChunk,
					},
					Content: stringChunk,
				})
			}

			return nil
		}),
		llms.WithTools(llmTools),
		llms.WithN(llmChat.N),
		llms.WithMaxTokens(llmChat.MaxTokens),
		llms.WithTemperature(llmChat.Temperature),
		llms.WithTopP(llmChat.TopP),
		llms.WithTopK(llmChat.TopK))
	return resp, err
}

func (s *Service) sendResponse(responseChan chan<- *schema.AssistantResponse, response schema.AssistantResponse) error {
	// 使用 select 来尝试写入数据，检测 channel 是否关闭
	select {
	case responseChan <- &response:
		// 正常写入
		return nil
	default:
		return io.EOF
	}
}

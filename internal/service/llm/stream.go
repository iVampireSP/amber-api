package llm

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/mitchellh/mapstructure"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/builtin_tool"
	"strconv"
)

// StreamChat 执行对话
func (s *Service) StreamChat(ctx context.Context, llmChat *schema.LLMChat, history []*entity.ChatMessage) error {
	// 不要从接收侧关闭 channel
	defer close(llmChat.ResponseChan)

	// 处理历史
	h, err := s.processHistory(ctx, llmChat, history)
	if err != nil {
		return err
	}

	var historyContent = h.MessageContent

	// 是否再次请求
	var requestAgain = true
	// 连续的函数调用次数
	var functionCallCount = 0
	//  全部的调用次数
	var totalFunctionCallCount = 0

	var tokenUsage = &schema.TokenUsage{}

	for {
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

		// Built-in 工具的忽略参数
		var without = &builtin_tool.WithoutOptions{
			Image:    llmChat.WithoutImage,
			Browsing: llmChat.WithoutBrowsing,
			Search:   llmChat.WithoutInternetSearch,
		}

		var llmTools = append(s.BuiltInTools.GetTools(without), llmChat.Tools...)

		resp, err := s.GenerateContent(ctx, llmChat, llmTools, historyContent)
		if err != nil {

			s.write(ctx, llmChat, &schema.AssistantResponse{
				State:   schema.StateFailed,
				Content: err.Error(),
			})
			return err
		}

		respChoice := resp.Choices[0]
		tokenUsage2 := s.getTokenUsage(respChoice)
		tokenUsage.CompletionTokens += tokenUsage2.CompletionTokens
		tokenUsage.PromptTokens += tokenUsage2.PromptTokens
		tokenUsage.TotalTokens += tokenUsage2.TotalTokens

		if respChoice.FuncCall != nil {
			// 检测到工具调用，标记为需要再次请求
			requestAgain = true

			// 将计数器加一次
			functionCallCount += 1

			// 这个代码本质上是针对 qwen 的，但是 OpenAI 不会有这样需要拼接的问题。
			// 拼接完整参数
			//var fullArgs = ""
			//for _, tc := range respChoice.ToolCalls {
			//	// 拼接参数
			//	fullArgs += tc.FunctionCall.Arguments
			//}

			// 处理 ToolCall
			for _, tc := range respChoice.ToolCalls {
				// 设置 AI 消息，里面加入 tool call
				assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, respChoice.Content)
				assistantResponse.Parts = append(assistantResponse.Parts, tc)
				historyContent = append(historyContent, assistantResponse)

				// 去除 fullArgs 的首尾 \n（一直检测）
				//for {
				//	if fullArgs[0] == '\n' {
				//		fullArgs = fullArgs[1:]
				//	} else if fullArgs[len(fullArgs)-1] == '\n' {
				//		fullArgs = fullArgs[:len(fullArgs)-1]
				//	} else {
				//		break
				//	}
				//}

				// 解析工具
				var functionCallArgs schema.FunctionCallArguments
				err = sonic.Unmarshal([]byte(tc.FunctionCall.Arguments), &functionCallArgs)
				if err != nil {
					return err
				}

				prefix, functionName := s.spiltFunctionName(tc.FunctionCall.Name)

				var toolCalling = &schema.AssistantResponse{
					State: schema.StateToolCalling,
					ToolCallMessage: &schema.ToolCallMessage{
						FunctionName: tc.FunctionCall.Name,
						Arguments:    functionCallArgs,
					},
					Internal: &schema.AssistantInternal{
						ToolCall:   &tc,
						ToolCallId: tc.ID,
					},
				}

				var toolRemoteResponse = &schema.ToolRemoteResponse{}
				var toolName = ""

				var selectedTool *entity.Tool

				// Built-in tools
				if prefix == builtin_tool.NAME {
					toolName = "Built-in"
				} else {
					// 转换 prefix
					toolId, err := strconv.Atoi(prefix)
					if err != nil {
						s.write(ctx, llmChat, &schema.AssistantResponse{
							// 这里改成 failed 会不会更好？
							State:   schema.StateToolFailed,
							Content: err.Error(),
							ToolResponseMessage: &schema.ToolResponseMessage{
								ToolName:     builtin_tool.NAME,
								FunctionName: tc.FunctionCall.Name,
								Content:      err.Error(),
							},
							Internal: &schema.AssistantInternal{
								ToolCall:   &tc,
								ToolCallId: tc.ID,
							},
							TokenUsage: tokenUsage,
						})
						return err
					}

					// 获取 Tool
					selectedTool, err = s.GetToolById(ctx, schema.EntityId(toolId))
					if err != nil {
						s.write(ctx, llmChat, &schema.AssistantResponse{
							// 这里改成 failed 会不会更好？
							State:   schema.StateToolFailed,
							Content: err.Error(),
							ToolResponseMessage: &schema.ToolResponseMessage{
								ToolName:     toolName,
								FunctionName: tc.FunctionCall.Name,
								Content:      err.Error(),
							},
							TokenUsage: tokenUsage,
							Internal: &schema.AssistantInternal{
								ToolCall:   &tc,
								ToolCallId: tc.ID,
							},
						})
						return err
					}

					toolName = selectedTool.Name
				}

				toolCalling.ToolCallMessage.ToolName = toolName
				// 发布工具调用
				s.write(ctx, llmChat, toolCalling)

				if prefix == builtin_tool.NAME {
					// 是 builtin，则调用内置函数
					var builtInToolRequest = &schema.CallBuiltInToolRequest{
						FunctionName: functionName,
						Args:         functionCallArgs,
					}

					s.Logger.Sugar.Infof("Calling Builtin function: %v, args: %v", functionName, functionCallArgs)
					builtInResponse, err := s.BuiltInTools.CallFunction(ctx, builtInToolRequest)
					if err != nil {
						// 也许内置函数不应该报 ToolFailed,不如直接 failed
						s.write(ctx, llmChat, &schema.AssistantResponse{
							State:      schema.StateFailed,
							Content:    err.Error(),
							TokenUsage: tokenUsage,
							//Internal: &schema.AssistantInternal{
							//	ToolCall:   &tc,
							//	ToolCallId: tc.ID,
							//},
						})
						return err
					}

					s.Logger.Sugar.Infof("Builtin response: %v", builtInResponse.Content)

					toolName = builtin_tool.NAME

					// mapstructure
					err = mapstructure.Decode(builtInResponse, toolRemoteResponse)
					if err != nil {
						return err
					}

					if builtInResponse.TokenUsage != nil {
						tokenUsage.PromptTokens += builtInResponse.TokenUsage.PromptTokens
						tokenUsage.CompletionTokens += builtInResponse.TokenUsage.CompletionTokens
						tokenUsage.TotalTokens += builtInResponse.TokenUsage.TotalTokens
					}

					//toolCalling.ToolCallMessage.ToolName = builtin_tool.NAME
				} else {
					s.Logger.Sugar.Infof("Calling Remote function: %v, args: %v", functionName, functionCallArgs)

					// 调用远程函数
					toolRemoteResponse, err = s.callRemoteFunction(selectedTool, llmChat, functionName, functionCallArgs)
					if err != nil {
						s.write(ctx, llmChat, &schema.AssistantResponse{
							State:   schema.StateToolFailed,
							Content: err.Error(),
							ToolResponseMessage: &schema.ToolResponseMessage{
								ToolName:     toolName,
								FunctionName: tc.FunctionCall.Name,
								Content:      err.Error(),
							},
							TokenUsage: tokenUsage,
						})
						return err
					}
				}

				// 如果是 builtin ，则不告知
				//if toolName != builtin_tool.NAME {
				//}

				// 算了，告知吧，好处理点
				s.write(ctx, llmChat, &schema.AssistantResponse{
					State: schema.StateToolResponse,
					ToolResponseMessage: &schema.ToolResponseMessage{
						ToolName:       toolName,
						FunctionName:   tc.FunctionCall.Name,
						Content:        toolRemoteResponse.Content,
						StopGeneration: toolRemoteResponse.StopGeneration,
						Append:         toolRemoteResponse.Append,
						Role:           toolRemoteResponse.Role,
						Text:           toolRemoteResponse.Text,
					},
					Internal: &schema.AssistantInternal{
						ToolCallId: tc.ID,
						ToolCall:   &tc,
					},
				})

				// End Built-in tools

				// ToolCall 处理完成，放入 History
				historyContent = append(historyContent, llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: tc.ID,
							Name:       tc.FunctionCall.Name,
							Content:    toolRemoteResponse.Content,
						},
					},
				})

				// 如果函数要求停止生成
				if toolRemoteResponse.StopGeneration {
					requestAgain = false
				}

			}

			// ToolCall 处理完成
		} else {
			// 不是工具调用，不再进行新的一轮请求，然后清除计数器
			requestAgain = false
			totalFunctionCallCount += functionCallCount
			functionCallCount = 0
		}

		s.write(ctx, llmChat, &schema.AssistantResponse{
			State: schema.StateDone,
		})

		s.tokenUsage.IncrMonthTokenUsage(ctx, tokenUsage.TotalTokens)
		if totalFunctionCallCount > 0 {
			s.tokenUsage.IncrMonthToolCallTimes(ctx, totalFunctionCallCount)
		}

		historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))

		//fmt.Println("本轮历史", historyContent)
	}

	s.write(ctx, llmChat, &schema.AssistantResponse{
		State:      schema.StateFinished,
		TokenUsage: tokenUsage,
	})

	return nil
}

func (s *Service) write(ctx context.Context, llmChat *schema.LLMChat, r *schema.AssistantResponse) {
	llmChat.ResponseChan <- r

	go s.event(ctx, llmChat.UserPublicInfo, r)
}

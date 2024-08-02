package llm

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/entity"
	"strings"
)

type Callback func(string)

// StreamChat 执行对话
func (s *Service) StreamChat(responseChan chan *AssistantResponse, history []*entity.ChatMessage, tools ...llms.Tool) error {
	var historyContent []llms.MessageContent

	for _, h := range history {
		switch h.Role {
		case entity.RoleHuman:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeHuman, h.Content))
		case entity.RoleAssistant:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeAI, h.Content))
		case entity.RoleSystem:
			historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, h.Content))
		}
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

					responseChan <- &AssistantResponse{
						State: StateChunk,
						ChunkMessage: &ChunkMessage{
							Content: stringChunk,
						},
						Content: stringChunk,
					}

					fullResponse = append(fullResponse, chunk)
				}

				return nil
			}),
			llms.WithTools(tools))
		if err != nil {
			panic(err)
		}

		respChoice := resp.Choices[0]

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

			//fmt.Println("最终参数", fullArgs)

			// 解析工具
			var functionCallArgs FunctionCallArgs
			err = sonic.Unmarshal([]byte(fullArgs), &functionCallArgs)
			if err != nil {
				panic(err)
			}

			responseChan <- &AssistantResponse{
				State: StateToolCalling,
				ToolCallMessage: &ToolCallMessage{
					Name: respChoice.FuncCall.Name,
					Args: functionCallArgs,
				},
			}

			var fakeToolResponseContent = "天气晴天，气温 25°C"

			functionName := s.spiltFunctionName(respChoice.FuncCall.Name)
			switch functionName {
			case "get_current_weather":

				historyContent = append(historyContent, llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: respChoice.ToolCalls[0].ID,
							Name:       respChoice.FuncCall.Name,
							Content:    fakeToolResponseContent,
						},
					},
				})

			}

			responseChan <- &AssistantResponse{
				State: StateToolResponse,
				ToolResponseMessage: &ToolResponseMessage{
					Name:    respChoice.FuncCall.Name,
					Content: fakeToolResponseContent,
				},
			}

		} else {
			requestAgain = false
		}

		historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))

		//fmt.Println("本轮历史", historyContent)
	}

	responseChan <- &AssistantResponse{
		State: StateDone,
	}

	//close(responseChan)
	return nil
}

func (s *Service) spiltFunctionName(functionName string) string {
	// 根据 _ 分割
	var functionNames = strings.Split(functionName, "_")

	// 从第 1 个开始取到最后一个
	var toolName = strings.Join(functionNames[1:], "_")

	//fmt.Println("解析的工具名称：" + toolName)

	return toolName
}

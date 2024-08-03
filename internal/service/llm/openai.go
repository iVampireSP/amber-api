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
	"strconv"
	"strings"
)

// StreamChat 执行对话
func (s *Service) StreamChat(responseChan chan *AssistantResponse, systemPrompt string, history []*entity.ChatMessage, user *schema.UserTokenInfo, tools ...llms.Tool) error {
	var historyContent []llms.MessageContent

	if systemPrompt != "" {
		historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt))
	}

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

			tool, functionName, err := s.spiltFunctionName(respChoice.FuncCall.Name)
			if err != nil {
				responseChan <- &AssistantResponse{
					State:   StateFailed,
					Content: err.Error(),
				}
			}

			remoteFunctionResponse, err := s.callRemoteFunction(tool, user, functionName, functionCallArgs)
			if err != nil {
				responseChan <- &AssistantResponse{
					State:   StateToolFailed,
					Content: err.Error(),
					ToolResponseMessage: &ToolResponseMessage{
						Name:    respChoice.FuncCall.Name,
						Content: err.Error(),
					},
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

				responseChan <- &AssistantResponse{
					State: StateToolResponse,
					ToolResponseMessage: &ToolResponseMessage{
						Name:    respChoice.FuncCall.Name,
						Content: remoteFunctionResponse.Content,
					},
				}
			}

		} else {
			requestAgain = false
		}

		historyContent = append(historyContent, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))

		//fmt.Println("本轮历史", historyContent)
	}

	responseChan <- &AssistantResponse{
		State: StateFinished,
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

func (s *Service) callRemoteFunction(tool *entity.Tool, user *schema.UserTokenInfo, functionName string, args FunctionCallArgs) (*schema.ToolRemoteResponse, error) {
	var callbackUrl = tool.Data.CallbackUrl

	argsJson, err := args.JSON()
	if err != nil {
		return nil, err
	}

	var userPublicInfo = &schema.UserPublicInfo{
		Name: user.Name,
		Id:   user.Sub,
	}

	var toolRequest = &schema.ToolRemoteRequest{
		FunctionName: functionName,
		Parameters:   string(argsJson),
		ApiKey:       tool.ApiKey,
		User:         userPublicInfo,
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

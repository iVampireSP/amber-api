package main

import (
	"context"
	"fmt"
	"rag-new/internal/base"
	"rag-new/internal/entity"
	"rag-new/internal/service/llm"

	"github.com/bytedance/sonic"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func init() {
	rootCmd.AddCommand(llmCmd)
}

var llmCmd = &cobra.Command{
	Use:  "llm",
	Long: "测试 LLM",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := CreateApp()
		if err != nil {
			panic(err)
			return
		}

		historyCallLLM(app)

		//llm, err := openai.New(
		//	openai.WithToken(app.Config.OpenAI.ApiKey),
		//	openai.WithBaseURL(app.Config.OpenAI.BaseUrl),
		//	openai.WithModel(app.Config.OpenAI.Model),
		//)

		//if err != nil {
		//	log.Fatal(err)
		//}
		//callLLM(llm)

		//TestCall(app)

	},
}

type ChunkMessage struct {
	Type     string        `json:"type"`
	Function FunctionChunk `json:"function"`
}

type FunctionChunk struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type FunctionCallArgs map[string]interface{}

func historyCallLLM(app *base.Application) {
	// fake history
	histories := []*entity.ChatMessage{
		{
			Role:    entity.RoleHuman,
			Content: "北京天气如何",
		},
	}

	var responseChan = make(chan *llm.AssistantResponse)

	go func() {
		fmt.Print(">> ")

		for {
			select {
			case msg := <-responseChan:
				//fmt.Println("接收到", msg.Content)
				fmt.Print(msg.Content)

				if msg.State == llm.StateDone {
					fmt.Println("")
					// 关闭 chan
					close(responseChan)

					break
				}
			}
		}

	}()

	err := app.Service.LLM.StreamChat(responseChan, histories)
	if err != nil {
		return
	}

	//

}

func callLLM(llm *openai.LLM) {

	var histories []llms.MessageContent
	var requestAgain = true

	for {
		var question string
		fmt.Print(">> ")
		n, err := fmt.Scanln(&question)
		if n == 0 || err != nil {
			question = "北京天气怎么样"
			fmt.Print(question + "\n")
			//continue
		}

		requestAgain = true

		histories = append(histories, llms.TextParts(llms.ChatMessageTypeHuman, question))

		var fullResponse [][]byte

		ctx := context.Background()

		for {
			fmt.Println("再次请求吗?" + fmt.Sprint(requestAgain))
			if !requestAgain {
				break
			}

			requestAgain = false
			var isToolCall = false

			fmt.Println("对话历史", histories)

			resp, err := llm.GenerateContent(ctx,
				histories,
				llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
					// 检测长度
					if len(chunk) == 0 {
						return nil
					}

					fmt.Printf("Received chunk: %s\n", chunk)

					fullResponse = append(fullResponse, chunk)

					return nil
				}),
				llms.WithTools(tools))
			if err != nil {
				panic(err)
			}

			respchoice := resp.Choices[0]

			if respchoice.FuncCall != nil {
				fmt.Println("FunCall 检测到工具调用")
				isToolCall = true
				requestAgain = true
			}

			if isToolCall {
				fmt.Printf("正在调用: %v\n", respchoice.FuncCall.Name)

				var fullArgs = ""

				assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, respchoice.Content)
				for _, tc := range respchoice.ToolCalls {
					// 拼接参数
					fullArgs += tc.FunctionCall.Arguments
				}

				var firstToolCall = respchoice.ToolCalls[0]
				firstToolCall.FunctionCall.Arguments = fullArgs
				assistantResponse.Parts = append(assistantResponse.Parts, firstToolCall)

				histories = append(histories, assistantResponse)

				fmt.Println("最终参数", fullArgs)

				// 解析工具
				var functionCallArgs FunctionCallArgs
				err = sonic.Unmarshal([]byte(fullArgs), &functionCallArgs)
				if err != nil {
					panic(err)
				}

				switch respchoice.FuncCall.Name {
				case "getCurrentWeather":
					weatherInfo, err := getCurrentWeather(functionCallArgs["location"].(string), functionCallArgs["unit"].(string))
					if err != nil {
						panic(err)
					}

					histories = append(histories, llms.MessageContent{
						Role: llms.ChatMessageTypeTool,
						Parts: []llms.ContentPart{
							llms.ToolCallResponse{
								ToolCallID: respchoice.ToolCalls[0].ID,
								Name:       respchoice.FuncCall.Name,
								Content:    weatherInfo,
							},
						},
					})

				}

			} else {
				requestAgain = false
			}

			histories = append(histories, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))

			fmt.Println("本轮历史", histories)
		}

	}
}

func getCurrentWeather(location string, unit string) (string, error) {
	//weatherInfo := map[string]interface{}{
	//	"location":    location,
	//	"temperature": "72",
	//	"unit":        unit,
	//	"forecast":    []string{"sunny", "windy"},
	//}
	//b, err := json.Marshal(weatherInfo)
	//if err != nil {
	//	return "", err
	//}
	//return string(b), nil

	return "晴天", nil
}

// json.RawMessage(`{"type": "object", "properties": {"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}}, "required": ["location"]}`),

var tools = []llms.Tool{
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "getCurrentWeather",
			Description: "Get the current weather in a given location",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"rationale": {
						Type:        jsonschema.String,
						Description: "The rationale for choosing this function call with these parameters",
					},
					"location": {
						Type:        jsonschema.String,
						Description: "The city and state, e.g. San Francisco, CA",
					},
					"unit": {
						Type: jsonschema.String,
						Enum: []string{"celsius", "fahrenheit"},
					},
				},
				Required: []string{"rationale", "location"},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "getTomorrowWeather",
			Description: "Get the predicted weather in a given location",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"rationale": {
						Type:        jsonschema.String,
						Description: "The rationale for choosing this function call with these parameters",
					},
					"location": {
						Type:        jsonschema.String,
						Description: "The city and state, e.g. San Francisco, CA",
					},
					"unit": {
						Type: jsonschema.String,
						Enum: []string{"celsius", "fahrenheit"},
					},
				},
				Required: []string{"rationale", "location"},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "getSuggestedPrompts",
			Description: "Given the user's input prompt suggest some related prompts",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"rationale": {
						Type:        jsonschema.String,
						Description: "The rationale for choosing this function call with these parameters",
					},
					"suggestions": {
						Type: jsonschema.Array,
						Items: &jsonschema.Definition{
							Type:        jsonschema.String,
							Description: "A suggested prompt",
						},
					},
				},
				Required: []string{"rationale", "suggestions"},
			},
		},
	},
}

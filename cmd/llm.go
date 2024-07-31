package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"log"
	"rag-new/internal/base"
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

		//TestLLM(app)
		TestCall(app)

	},
}

func TestLLM(app *base.Application) {
	ctx := context.Background()

	llm, err := openai.New(
		openai.WithToken(app.Config.OpenAI.ApiKey),
		openai.WithBaseURL(app.Config.OpenAI.BaseUrl),
		openai.WithModel("qwen-max"),
	)
	if err != nil {
		app.Logger.Sugar.Fatal(err)
	}
	prompt := "What would be a good company name for a company that makes colorful socks?"
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
}

func TestCall(app *base.Application) {
	llm, err := openai.New(
		openai.WithToken(app.Config.OpenAI.ApiKey),
		openai.WithBaseURL(app.Config.OpenAI.BaseUrl),
		openai.WithModel("qwen-max"),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	resp, err := llm.GenerateContent(ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, "北京天气怎么样"),
		},
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Printf("Received chunk: %s\n", chunk)
			return nil
		}),
		llms.WithTools(tools))
	if err != nil {
		log.Fatal(err)
	}

	choice1 := resp.Choices[0]
	if choice1.FuncCall != nil {
		fmt.Printf("Function call: %v\n", choice1.FuncCall)
	}
}

func getCurrentWeather(location string, unit string) (string, error) {
	weatherInfo := map[string]interface{}{
		"location":    location,
		"temperature": "72",
		"unit":        unit,
		"forecast":    []string{"sunny", "windy"},
	}
	b, err := json.Marshal(weatherInfo)
	if err != nil {
		return "", err
	}
	return string(b), nil
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

package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
)

func init() {
	RootCmd.AddCommand(toolCmd)
}

var toolCmd = &cobra.Command{
	Use:  "tool",
	Long: "测试 LLM Tool",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := CreateApp()
		if err != nil {
			panic(err)
			return
		}

		var ctx = context.Background()

		assistant, err := app.Service.Assistant.GetAssistant(ctx, 1)
		if err != nil {
			panic(err)
		}

		tool, err := app.Service.Assistant.ToLLMTool(ctx, assistant)
		if err != nil {
			panic(err)
			return
		}

		fmt.Println(tool)

		//var defaultUser schema.UserId = 1
		//
		//toolEntity, err := app.Service.Tool.ListToolFromUserId(ctx, defaultUser)
		//if err != nil {
		//	panic(err)
		//}
		//
		//var toolList []llms.Tool
		//
		//// 转换格式
		//for _, v := range toolEntity {
		//	var fnData = v.Data
		//	var llmTool = llms.Tool{
		//		Type: "function",
		//	}
		//
		//	for _, v := range fnData.ToolFunctions {
		//		for _, j := range v.Function {
		//			llmTool.Function = &llms.FunctionDefinition{
		//				Name:        j.Name,
		//				Description: j.Description,
		//				Parameters:  j.Parameters,
		//			}
		//		}
		//	}
		//	toolList = append(toolList, llmTool)
		//
		//	//toolList = append(toolList, llms.Tool{
		//	//	Type: "function",
		//	//	Function: &llms.FunctionDefinition{
		//	//		Name:        fnData.Name,
		//	//		Description: fnData.Description,
		//	//		Parameters:  fnData.ToolFunctions[0].Function[0].Parameters,
		//	//	},
		//	//})
		//}
		//
		//fmt.Println(toolList[0].Function.Name)

	},
}

var tools2 = []llms.Tool{
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

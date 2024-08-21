package builtin_tool

import (
	"context"
	"errors"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
)

var NAME = "builtin"

var tools = []llms.Tool{
	//{
	//	Type: "function",
	//	Function: &llms.FunctionDefinition{
	//		Name:        prefix("now"),
	//		Description: "Get the server time using server's timezone(not users)",
	//		Parameters: jsonschema.Definition{
	//			Type: jsonschema.Object,
	//			Properties: map[string]jsonschema.Definition{
	//				"rationale": {
	//					Type:        jsonschema.String,
	//					Description: "The rationale for choosing this function call with these parameters",
	//				},
	//			},
	//		},
	//	},
	//},
}

func prefix(name string) string {
	return NAME + "_" + name
}

func (s *Service) GetTools() []llms.Tool {
	return tools
}

func (s *Service) CallFunction(ctx context.Context, functionName string, args schema.FunctionCallArguments) (*schema.ToolRemoteResponse, error) {
	var response = &schema.ToolRemoteResponse{}
	switch functionName {
	case "now":
		response.Content = s.GetCurrentTime()
	default:
		return nil, errors.New("function not found")
	}

	return response, nil
}

// Exists 这里拿到的 functionName 是不带前缀的，如果 withPrefix 为 true，则带前缀传入并判断
func (*Service) Exists(functionName string, withPrefix bool) bool {
	for _, tool := range tools {
		if withPrefix {
			if NAME+"_"+tool.Function.Name == functionName {
				return true
			}
		} else {
			if tool.Function.Name == functionName {
				return true
			}
		}
	}
	return false
}

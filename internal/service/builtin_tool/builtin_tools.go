package builtin_tool

import (
	"context"
	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
)

var PREFIX = "builtin_"
var NAME = "builtin"

var tools = []llms.Tool{
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "now",
			Description: "Get the server time using server's timezone(not users)",
			Parameters:  jsonschema.Definition{},
		},
	},
}

func (s *Service) GetTools() []llms.Tool {
	// 为每个 function name 加上 prefix
	for i := range tools {
		tools[i].Function.Name = PREFIX + tools[i].Function.Name
	}

	return tools
}

func (s *Service) CallFunction(ctx context.Context, functionName string, args schema.FunctionCallArguments) (*schema.ToolRemoteResponse, error) {
	var response = &schema.ToolRemoteResponse{}
	switch functionName {
	case "now":
		response.Content = s.GetCurrentTime()
	}

	return response, nil
}

// Exists 这里拿到的 functionName 是不带前缀的，如果 withPrefix 为 true，则带前缀传入并判断
func (*Service) Exists(functionName string, withPrefix bool) bool {
	for _, tool := range tools {
		if withPrefix {
			if PREFIX+tool.Function.Name == functionName {
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

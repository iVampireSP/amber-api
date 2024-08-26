package builtin_tool

import (
	"context"
	"errors"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
)

var NAME = "builtin"

func prefix(name string) string {
	return NAME + "_" + name
}

func (s *Service) GetTools() []llms.Tool {
	return tools
}

func (s *Service) CallFunction(ctx context.Context, functionName string, args schema.FunctionCallArguments) (*schema.ToolRemoteResponse, error) {
	var response = &schema.ToolRemoteResponse{}
	var err error = nil
	var textResponse string

	switch functionName {
	case "now":
		response.Content = s.GetCurrentTime()
	case "describe_image":
		textResponse, err = s.DescribeImage(ctx, args)

		if err != nil {
			response.Success = false
			response.StopGeneration = true
		} else {
			response.Content = textResponse
		}

	default:
		return nil, errors.New("function not found")
	}

	return response, err
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

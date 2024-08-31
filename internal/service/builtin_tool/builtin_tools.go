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

type WithoutOptions struct {
	Image bool
}

func (s *Service) GetTools(without *WithoutOptions) []llms.Tool {
	if without == nil {
		return tools
	}

	var t []llms.Tool
	for _, v := range tools {
		if v.Function.Name == prefix("describe_image") {
			if without.Image {
				continue
			} else {
				t = append(t, v)
			}
		}
	}

	return t
}

func (s *Service) CallFunction(ctx context.Context, req *schema.CallBuiltInToolRequest) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}
	var err error = nil

	switch req.FunctionName {
	case "now":
		response.Content = s.GetCurrentTime()
	case "describe_image":
		response, err = s.DescribeImage(ctx, req.Args)

		//if err != nil {
		//	response.Success = false
		//	response.StopGeneration = false
		//}
	case "download_file":
		response, err = s.DownloadFile(ctx, req.Args)
	default:
		return nil, errors.New("function not found")
	}

	if err != nil {
		s.logger.Sugar.Error("Built-in failed: " + err.Error())
		// reset response
		response = &schema.CallBuiltInResponse{}
		response.Content = "Built-in tool Error: " + err.Error()
		return response, nil
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

package builtin_tool

import (
	"context"
	"errors"
	"rag-new/internal/schema"
)

var NAME = "builtin"

func prefix(name string) string {
	return NAME + "_" + name
}

func (s *Service) CallFunction(ctx context.Context, req *schema.CallBuiltInToolRequest) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}
	var err error = nil

	switch req.FunctionName {
	case "describe_image":
		response, err = s.DescribeImage(ctx, req.Args)
	case "generate_image":
		response, err = s.GenerateImage(ctx, req.Args)
	case "calculator":
		response, err = s.Calculator(ctx, req.Args)
	case "compare":
		response, err = s.Compare(ctx, req.Args)
	case "search_web":
		response, err = s.SearchWeb(ctx, req.Args)
	case "get_url_content":
		response, err = s.ReadUrl(ctx, req.Args)
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

//// Exists 这里拿到的 functionName 是不带前缀的，如果 withPrefix 为 true，则带前缀传入并判断
//func (*Service) Exists(functionName string, withPrefix bool) bool {
//	for _, tool := range tools {
//		if withPrefix {
//			if NAME+"_"+tool.Function.Name == functionName {
//				return true
//			}
//		} else {
//			if tool.Function.Name == functionName {
//				return true
//			}
//		}
//	}
//	return false
//}

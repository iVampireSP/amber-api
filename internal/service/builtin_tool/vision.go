package builtin_tool

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
)

type describeImageParams struct {
	Question string          `json:"question"`
	ImageId  schema.EntityId `json:"image_id" mapstructure:"image_id"`
}

func (s *Service) DescribeImage(ctx context.Context, args schema.FunctionCallArguments) (string, error) {
	// TODO: 计算 Token 消耗
	var params describeImageParams
	err := args.Unmarshal(&params)
	if err != nil {
		return "", err
	}

	// 文件必须存在
	exists, err := s.fileService.ExistsFileById(ctx, params.ImageId)
	if err != nil {
		return "此时无法获取文件是否存在", nil
	}
	if !exists {
		return "文件不存在", nil
	}

	// 获取文件
	file, err := s.fileService.GetFileById(ctx, params.ImageId)
	if err != nil {
		return "此时无法获取文件", nil
	}

	// URL
	fileUrl, err := s.fileService.GetImageUrl(file)
	if err != nil {
		return "此时无法获取文件 Url", nil
	}

	var describeImageHistory = []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Provide a brief response in the user's language"),
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.ImageURLWithDetailPart(fileUrl, "auto"),
				llms.TextPart(params.Question),
			},
		},
	}

	resp, err := s.OpenAI.GenerateContent(ctx, describeImageHistory)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Content, nil

}

package builtin_tool

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
	"strconv"
)

type describeImageParams struct {
	Question string `json:"question"`
	ImageId  string `json:"image_id" mapstructure:"image_id"`
}

func (s *Service) DescribeImage(ctx context.Context, args schema.FunctionCallArguments) (string, error) {
	var params describeImageParams
	err := args.Unmarshal(&params)
	if err != nil {
		return "", err
	}

	imageIdInt, err := strconv.ParseInt(params.ImageId, 10, 64)

	// 文件必须存在
	exists, err := s.fileService.ExistsFileById(ctx, schema.EntityId(imageIdInt))
	if err != nil {
		return "此时无法获取文件是否存在", nil
	}
	if !exists {
		return "文件不存在", nil
	}

	// 获取文件
	file, err := s.fileService.GetFileById(ctx, schema.EntityId(imageIdInt))
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
				llms.ImageURLPart(fileUrl),
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

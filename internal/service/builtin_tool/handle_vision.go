package builtin_tool

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
)

type describeImageParams struct {
	Prompt string `json:"prompt"`
	Url    string `json:"url"`
	Hash   string `json:"hash" mapstructure:"hash"`
}

func (s *Service) DescribeImage(ctx context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}
	var params describeImageParams
	err := args.Unmarshal(&params)
	if err != nil {
		return response, err
	}

	if params.Url == "" && params.Hash == "" {
		response.Content = "没有图片 URL 或者 Hash"
		return response, nil
	}

	var file = &entity.File{}

	if params.Url != "" {
		file, err = s.fileService.CreateFileFromUrl(ctx, params.Url)
		if err != nil {
			return response, err
		}
	} else {
		// 文件必须存在
		exists, err := s.fileService.ExistsFileByFileHash(ctx, params.Hash)
		if err != nil {
			return response, err
		}
		if !exists {
			response.Content = "文件不存在"

			return response, nil
		}

		// 获取文件
		file, err = s.fileService.GetFileByFileHash(ctx, params.Hash)
		if err != nil {
			response.Content = "此时无法获取文件"

			return response, nil
		}

	}

	// 如果不是图片
	if !file.IsImage() {
		response.Content = "文件不是图片"

		return response, nil
	}

	// URL
	fileUrl, err := s.fileService.GetImageUrl(file)
	if err != nil {
		response.Content = "此时无法获取文件 URL"

		return response, nil
	}

	var describeImageHistory = []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Provide a brief response in the user's language"),
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.ImageURLWithDetailPart(fileUrl, "auto"),
				llms.TextPart(params.Prompt),
			},
		},
	}

	resp, err := s.LLM.GenerateContent(ctx, describeImageHistory)
	if err != nil {
		return response, err
	}

	var tokenUsage = s.getTokenUsage(resp.Choices[0])
	response.Content = resp.Choices[0].Content
	response.TokenUsage = tokenUsage
	return response, nil
}

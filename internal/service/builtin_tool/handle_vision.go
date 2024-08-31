package builtin_tool

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"strings"
)

type describeImageParams struct {
	Question string          `json:"question"`
	Url      string          `json:"url"`
	FileId   schema.EntityId `json:"file_id" mapstructure:"file_id"`
}

func (s *Service) DescribeImage(ctx context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}
	var params describeImageParams
	err := args.Unmarshal(&params)
	if err != nil {
		return response, err
	}

	if params.Url == "" && params.FileId == 0 {
		response.Content = "请提供图片 URL 或者文件 ID"
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
		exists, err := s.fileService.ExistsFileById(ctx, params.FileId)
		if err != nil {
			return response, err
		}
		if !exists {
			response.Content = "文件不存在"

			return response, nil
		}

		// 获取文件
		file, err = s.fileService.GetFileById(ctx, params.FileId)
		if err != nil {
			response.Content = "此时无法获取文件"

			return response, nil
		}

	}

	// 如果 mimetype 不是 image/ 开头
	if !strings.HasPrefix(file.MimeType, "image/") {
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
				llms.TextPart(params.Question),
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

func (s *Service) getTokenUsage(respChoice *llms.ContentChoice) *schema.TokenUsage {
	var tokenUsage = &schema.TokenUsage{}

	// 如果 respChoice.GenerationInfo 中有 prompt_tokens
	if respChoice.GenerationInfo["PromptTokens"] != nil {
		tokenUsage.PromptTokens = respChoice.GenerationInfo["PromptTokens"].(int)
	} else {
		tokenUsage.PromptTokens = 0
	}

	// 如果 respChoice.GenerationInfo 中有 completion_tokens
	if respChoice.GenerationInfo["CompletionTokens"] != nil {
		tokenUsage.CompletionTokens = respChoice.GenerationInfo["CompletionTokens"].(int)
	} else {
		tokenUsage.CompletionTokens = 0
	}

	// 如果 respChoice.GenerationInfo 中有 total_tokens
	if respChoice.GenerationInfo["TotalTokens"] != nil {
		tokenUsage.TotalTokens = respChoice.GenerationInfo["TotalTokens"].(int)
	} else {
		tokenUsage.TotalTokens = 0
	}

	return tokenUsage
}

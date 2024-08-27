package builtin_tool

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
	"strings"
)

type describeImageParams struct {
	Question string          `json:"question"`
	ImageId  schema.EntityId `json:"image_id" mapstructure:"image_id"`
}

func (s *Service) DescribeImage(ctx context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	// TODO: 计算 Token 消耗
	var response = &schema.CallBuiltInResponse{}
	var params describeImageParams
	err := args.Unmarshal(&params)
	if err != nil {
		return nil, err
	}

	// 文件必须存在
	exists, err := s.fileService.ExistsFileById(ctx, params.ImageId)
	if err != nil {
		return response, err
	}
	if !exists {
		response.Content = "文件不存在"

		return response, nil
	}

	// 获取文件
	file, err := s.fileService.GetFileById(ctx, params.ImageId)
	if err != nil {
		response.Content = "此时无法获取文件"

		return response, nil
	}

	// 如果 mimetype 不是 image/ 开头
	if !strings.HasPrefix(file.MimeType, "image/") {
		response.Content = "文件不是图片"

		return nil, nil
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

	resp, err := s.OpenAI.GenerateContent(ctx, describeImageHistory)
	if err != nil {
		return nil, err
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

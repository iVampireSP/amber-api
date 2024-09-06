package builtin_tool

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"rag-new/internal/schema"
	"slices"
)

var dallEAllowedSizes = []string{"1024x1792", "1792x1024"}

type generateImageParams struct {
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
}

func (s *Service) GenerateImage(ctx context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}
	var params generateImageParams
	err := args.Unmarshal(&params)
	if err != nil {
		return response, err
	}

	// 如果 size 不在合法范围内
	if !slices.Contains(dallEAllowedSizes, params.Size) {
		return response, fmt.Errorf("size must be one of %v", dallEAllowedSizes)
	}

	reqBase64 := openai.ImageRequest{
		Model:          s.config.OpenAI.DallEModel,
		Prompt:         params.Prompt,
		Size:           params.Size,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := s.OpenAI.CreateImage(ctx, reqBase64)
	if err != nil {
		return response, err
	}
	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		return response, err
	}

	// 读取为 io.ReadSeeker
	r := bytes.NewReader(imgBytes)

	fileEntity, err := s.fileService.CreateFile(ctx, r, true)

	url, err := s.fileService.GetImageUrl(fileEntity)
	if err != nil {
		return response, err
	}

	response.Content = url

	return response, nil
}

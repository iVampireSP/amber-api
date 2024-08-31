package builtin_tool

import (
	openai2 "github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms/openai"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/service/file"
)

type Service struct {
	LLM         *openai.LLM
	config      *conf.Config
	logger      *logger.Logger
	fileService *file.Service
	OpenAI      *openai2.Client
}

func NewService(config *conf.Config, logger *logger.Logger, fileService *file.Service, openAI *openai2.Client) *Service {
	llm, err := openai.New(
		openai.WithToken(config.OpenAI.ApiKey),
		openai.WithBaseURL(config.OpenAI.BaseUrl),
		openai.WithModel(config.OpenAI.VisionModel),
	)

	if err != nil {
		panic(err)
	}

	return &Service{
		llm, config, logger, fileService, openAI,
	}
}

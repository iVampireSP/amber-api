package embedding

import (
	"github.com/tmc/langchaingo/llms/openai"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/dao"
)

type Service struct {
	OpenAI *openai.LLM
	Logger *logger.Logger
	config *conf.Config
	dao    *dao.Query
}

func NewService(config *conf.Config, logger *logger.Logger, dao *dao.Query) *Service {
	llm, err := openai.New(
		openai.WithToken(config.OpenAI.ApiKey),
		openai.WithBaseURL(config.OpenAI.InternalBaseUrl),
		openai.WithEmbeddingModel(config.OpenAI.EmbeddingModel),
	)

	if err != nil {
		panic(err)
	}

	return &Service{llm, logger, config, dao}
}

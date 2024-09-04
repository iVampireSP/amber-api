package memory

import (
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/tmc/langchaingo/llms/openai"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/dao"
	"rag-new/internal/service/embedding"
	"rag-new/internal/service/stream"
)

type Service struct {
	OpenAI    *openai.LLM
	Logger    *logger.Logger
	Embedding *embedding.Service
	Milvus    client.Client
	Stream    *stream.Service
	config    *conf.Config
	dao       *dao.Query
}

func NewMemory(config *conf.Config, logger *logger.Logger, embedding *embedding.Service, milvus client.Client, dao *dao.Query, stream *stream.Service) *Service {
	llm, err := openai.New(
		openai.WithToken(config.OpenAI.ApiKey),
		openai.WithBaseURL(config.OpenAI.BaseUrl),
	)

	if err != nil {
		panic(err)
	}

	return &Service{llm, logger, embedding, milvus, stream, config, dao}
}

package llm

import (
	"github.com/tmc/langchaingo/llms/openai"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/builtin_tool"
	"rag-new/internal/service/tool"
)

type Service struct {
	OpenAI           *openai.LLM
	Logger           *logger.Logger
	AssistantService *assistant.Service
	ToolService      *tool.Service
	BuiltInTools     *builtin_tool.Service
	config           *conf.Config
}

func NewLLM(config *conf.Config, logger *logger.Logger, assistantService *assistant.Service, toolService *tool.Service, builtinTools *builtin_tool.Service) *Service {
	llm, err := openai.New(
		openai.WithToken(config.OpenAI.ApiKey),
		openai.WithBaseURL(config.OpenAI.BaseUrl),
		openai.WithModel(config.OpenAI.Model),
	)

	if err != nil {
		panic(err)
	}

	return &Service{llm, logger, assistantService, toolService, builtinTools, config}
}

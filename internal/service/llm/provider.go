package llm

import (
	"github.com/tmc/langchaingo/llms/openai"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/dao"
	"rag-new/internal/message"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/builtin_tool"
	"rag-new/internal/service/chat"
	"rag-new/internal/service/file"
	"rag-new/internal/service/stream"
	"rag-new/internal/service/token_usage"
	"rag-new/internal/service/tool"
)

type Service struct {
	OpenAI           *openai.LLM
	Logger           *logger.Logger
	AssistantService *assistant.Service
	ToolService      *tool.Service
	BuiltInTools     *builtin_tool.Service
	FileService      *file.Service
	streamService    *stream.Service
	chatService      *chat.Service
	// 也许要把所有的名字改成小写，所以就从接下来的开始改成小写
	// 然后再慢慢改原来的吧
	message    *message.Message
	config     *conf.Config
	dao        *dao.Query
	tokenUsage *token_usage.Service
}

func NewLLM(config *conf.Config,
	logger *logger.Logger,
	assistantService *assistant.Service,
	toolService *tool.Service,
	builtinTools *builtin_tool.Service,
	fileService *file.Service,
	streamService *stream.Service,
	message *message.Message,
	dao *dao.Query,
	chat *chat.Service,
	tokenUsage *token_usage.Service,
) *Service {
	llm, err := openai.New(
		openai.WithToken(config.OpenAI.ApiKey),
		openai.WithBaseURL(config.OpenAI.BaseUrl),
	)

	if err != nil {
		panic(err)
	}

	return &Service{llm,
		logger,
		assistantService,
		toolService,
		builtinTools,
		fileService,
		streamService,
		chat,
		message,
		config,
		dao,
		tokenUsage,
	}
}

package service

import (
	"github.com/google/wire"
	"rag-new/internal/base/logger"
	"rag-new/internal/batch"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/builtin_tool"
	"rag-new/internal/service/chat"
	"rag-new/internal/service/chat_message"
	"rag-new/internal/service/embedding"
	"rag-new/internal/service/file"
	"rag-new/internal/service/jwks"
	"rag-new/internal/service/library"
	"rag-new/internal/service/llm"
	"rag-new/internal/service/memory"
	"rag-new/internal/service/message_block"
	"rag-new/internal/service/stream"
	"rag-new/internal/service/tool"
)

type Service struct {
	logger       *logger.Logger
	Jwks         *jwks.JWKS
	Auth         *auth.Service
	Tool         *tool.Service
	Assistant    *assistant.Service
	Chat         *chat.Service
	LLM          *llm.Service
	MessageBlock *message_block.Service
	ChatMessage  *chat_message.Service
	Batch        *batch.Batch
	BuiltinTool  *builtin_tool.Service
	File         *file.Service
	Stream       *stream.Service
	Library      *library.Service
	Embedding    *embedding.Service
}

var Provider = wire.NewSet(
	jwks.NewJWKS,
	auth.NewAuthService,
	embedding.NewService,
	memory.NewMemory,
	chat_message.NewService,
	chat.NewService,
	tool.NewService,
	assistant.NewService,
	builtin_tool.NewService,
	llm.NewLLM,
	message_block.NewService,
	file.NewService,
	stream.NewService,
	library.NewService,
	NewService,
)

func NewService(
	logger *logger.Logger,
	jwks *jwks.JWKS,
	auth *auth.Service,
	tool *tool.Service,
	assistant *assistant.Service,
	chat *chat.Service,
	llm *llm.Service,
	messageBlock *message_block.Service,
	chatMessage *chat_message.Service,
	builtinTool *builtin_tool.Service,
	batch *batch.Batch,
	file *file.Service,
	stream *stream.Service,
	library *library.Service,
	embedding *embedding.Service,
) *Service {
	return &Service{
		logger,
		jwks,
		auth,
		tool,
		assistant,
		chat,
		llm,
		messageBlock,
		chatMessage,
		batch,
		builtinTool,
		file,
		stream,
		library,
		embedding,
	}
}

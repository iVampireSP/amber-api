package service

import (
	"github.com/google/wire"
	"rag-new/internal/base/logger"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/chat"
	"rag-new/internal/service/jwks"
	"rag-new/internal/service/tool"
)

type Service struct {
	logger    *logger.Logger
	Jwks      *jwks.JWKS
	Auth      *auth.Service
	Tool      *tool.Service
	Assistant *assistant.Service
	Chat      *chat.Service
}

var Provider = wire.NewSet(
	jwks.NewJWKS,
	auth.NewAuthService,
	tool.NewService,
	assistant.NewService,
	chat.NewService,
	NewService,
)

func NewService(
	logger *logger.Logger,
	jwks *jwks.JWKS,
	auth *auth.Service,
	tool *tool.Service,
	assistant *assistant.Service,
	chat *chat.Service,
) *Service {
	return &Service{
		logger, jwks, auth, tool, assistant, chat,
	}
}

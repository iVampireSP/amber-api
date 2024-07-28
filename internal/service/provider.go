package service

import (
	"github.com/google/wire"
	"rag-new/internal/base/logger"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/jwks"
	"rag-new/internal/service/tool"
)

type Service struct {
	logger *logger.Logger
	Jwks   *jwks.JWKS
	Auth   *auth.Service
	Tool   *tool.Service
}

var Provider = wire.NewSet(
	jwks.NewJWKS,
	auth.NewAuthService,
	tool.NewService,
	NewService,
)

func NewService(
	jwks *jwks.JWKS,
	auth *auth.Service,
	logger *logger.Logger,
	tool *tool.Service,
) *Service {
	return &Service{
		Jwks:   jwks,
		Auth:   auth,
		logger: logger,
		Tool:   tool,
	}
}

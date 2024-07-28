package service

import (
	"github.com/google/wire"
	"rag-new/internal/base/logger"
	"rag-new/internal/service/auth"
	"rag-new/internal/service/jwks"
)

type Service struct {
	Jwks   *jwks.JWKS
	Auth   *auth.Service
	logger *logger.Logger
}

var Provider = wire.NewSet(
	jwks.NewJWKS,
	auth.NewAuthService,
	NewService,
)

func NewService(jwks *jwks.JWKS, auth *auth.Service, logger *logger.Logger) *Service {
	return &Service{
		Jwks:   jwks,
		Auth:   auth,
		logger: logger,
	}
}

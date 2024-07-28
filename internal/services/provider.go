package services

import (
	"github.com/google/wire"
	"rag-new/internal/logger"
	"rag-new/internal/services/auth"
	"rag-new/internal/services/jwks"
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

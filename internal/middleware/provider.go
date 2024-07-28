package middleware

import (
	"github.com/google/wire"
	"rag-new/internal/logger"
	"rag-new/internal/services/auth"
)

type Middleware struct {
	GinLogger   *GinLoggerMiddleware
	Auth        *AuthMiddleware
	AuthService *auth.Service
}

func NewMiddleware(logger *logger.Logger, authService *auth.Service) *Middleware {
	return &Middleware{
		GinLogger: NewGinLoggerMiddleware(logger.Logger.Desugar()),
		Auth:      NewAuthMiddleware(authService),
	}
}

var Provider = wire.NewSet(
	NewMiddleware,
)

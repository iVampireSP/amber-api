package middleware

import (
	"github.com/google/wire"
	"rag-new/internal/logger"
	"rag-new/internal/service/auth"
)

type Middleware struct {
	GinLogger   *GinLoggerMiddleware
	Auth        *AuthMiddleware
	AuthService *auth.Service
}

func NewMiddleware(logger *logger.Logger, authService *auth.Service) *Middleware {
	return &Middleware{
		GinLogger: NewGinLoggerMiddleware(logger.Logger),
		Auth:      NewAuthMiddleware(authService),
	}
}

var Provider = wire.NewSet(
	NewMiddleware,
)

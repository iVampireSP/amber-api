package middleware

import (
	"github.com/google/wire"
	"rag-new/internal/logger"
)

type Middleware struct {
	GinLogger *GinLoggerMiddleware
	Auth      *AuthMiddleware
}

func NewMiddleware(logger *logger.Logger) *Middleware {
	return &Middleware{
		GinLogger: NewGinLoggerMiddleware(logger.Logger.Desugar()),
		Auth:      NewAuthMiddleware(),
	}
}

var Provider = wire.NewSet(
	NewMiddleware,
)

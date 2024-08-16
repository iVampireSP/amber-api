package middleware

import (
	"github.com/google/wire"
	"rag-new/internal/base/logger"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
)

type Middleware struct {
	GinLogger              *GinLoggerMiddleware
	Auth                   *AuthMiddleware
	JSONResponse           *JSONResponseMiddleware
	AssistantTokenValidate *AssistantTokenValidateMiddleware
}

func NewMiddleware(logger *logger.Logger, authService *auth.Service, assistantService *assistant.Service) *Middleware {
	return &Middleware{
		GinLogger:              NewGinLoggerMiddleware(logger.Logger),
		Auth:                   NewAuthMiddleware(authService),
		JSONResponse:           NewJSONResponseMiddleware(),
		AssistantTokenValidate: NewAssistantTokenValidateMiddleware(assistantService),
	}
}

var Provider = wire.NewSet(
	NewMiddleware,
)

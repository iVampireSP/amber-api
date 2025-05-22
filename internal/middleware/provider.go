package middleware

import (
	"rag-new/internal/base/logger"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/auth"
)

type Middleware struct {
	GinLogger            *GinLoggerMiddleware
	Auth                 *AuthMiddleware
	JSONResponse         *JSONResponseMiddleware
	AssistantKeyValidate *AssistantKeyValidateMiddleware
}

func NewMiddleware(logger *logger.Logger, authService *auth.Service, assistantService *assistant.Service) *Middleware {
	return &Middleware{
		GinLogger:            NewGinLoggerMiddleware(logger.Logger),
		Auth:                 NewAuthMiddleware(authService),
		JSONResponse:         NewJSONResponseMiddleware(),
		AssistantKeyValidate: NewAssistantKeyValidateMiddleware(assistantService),
	}
}

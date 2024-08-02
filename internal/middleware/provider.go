package middleware

import (
	"github.com/google/wire"
	"rag-new/internal/base/logger"
	"rag-new/internal/service/auth"
)

type Middleware struct {
	GinLogger    *GinLoggerMiddleware
	Auth         *AuthMiddleware
	JSONResponse *JSONResponseMiddleware
	CORS         *CORSMiddleware
}

func NewMiddleware(logger *logger.Logger, authService *auth.Service) *Middleware {
	return &Middleware{
		GinLogger:    NewGinLoggerMiddleware(logger.Logger),
		Auth:         NewAuthMiddleware(authService),
		JSONResponse: NewJSONResponseMiddleware(),
		CORS:         NewCORSMiddleware(),
	}
}

var Provider = wire.NewSet(
	NewMiddleware,
)

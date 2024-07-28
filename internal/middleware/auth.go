package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
	"rag-new/pkg/consts"
)

type AuthMiddleware struct {
	authService *auth.Service
}

func NewAuthMiddleware(authService *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService,
	}
}

func (a AuthMiddleware) RequireJWTIDToken(c *gin.Context) {
	user, err := a.authService.GinMiddlewareAuth(schema.JWTIDToken, c)

	if err != nil {
		c.Abort()
		schema.ResponseError(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(consts.AuthMiddlewareKey, user)
	c.Next()
}

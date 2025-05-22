package middleware

import (
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/auth"
	"rag-new/pkg/consts"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *auth.Service
}

func NewAuthMiddleware(authService *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService,
	}
}

// RequireAuth 需要认证
func (a AuthMiddleware) RequireAuth(c *gin.Context) {
	user, err := a.authService.GinMiddlewareAuth(c)

	if err != nil {
		c.Abort()
		schema.NewResponse(c).Error(err).Status(http.StatusUnauthorized).Send()
		return
	}

	c.Set(consts.AuthMiddlewareKey, user)
	c.Next()
}

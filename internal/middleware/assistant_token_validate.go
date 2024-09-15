package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/assistant"
	"rag-new/pkg/consts"
	"strings"
)

type AssistantTokenValidateMiddleware struct {
	assistantService *assistant.Service
}

func NewAssistantTokenValidateMiddleware(assistantService *assistant.Service) *AssistantTokenValidateMiddleware {
	return &AssistantTokenValidateMiddleware{
		assistantService,
	}
}

func (a *AssistantTokenValidateMiddleware) AssistantTokenValidate(c *gin.Context) {
	var response = schema.NewResponse(c)
	authorization := c.Request.Header.Get(consts.AuthHeader)

	if authorization == "" {
		c.Abort()
		response.Error(consts.ErrBearerToken).Status(http.StatusUnauthorized).Send()
		return
	}

	authSplit := strings.Split(authorization, " ")
	if len(authSplit) != 2 {
		c.Abort()
		response.Error(consts.ErrBearerToken).Status(http.StatusUnauthorized).Send()
		return
	}

	if authSplit[0] != consts.AuthPrefix {
		c.Abort()
		response.Error(consts.ErrBearerToken).Status(http.StatusUnauthorized).Send()
		return
	}

	assistantToken := authSplit[1]
	// if start with sk-
	if strings.HasPrefix(assistantToken, "sk-") {
		assistantToken = assistantToken[3:]
	}

	if assistantToken == "" {
		c.Abort()
		response.Error(consts.ErrAssistantTokenNotFound).Status(http.StatusUnauthorized).Send()
		return
	}

	assistantEntity, err := a.assistantService.GetTokenBySecret(c, assistantToken)
	if assistantEntity == nil {
		c.Abort()
		response.Error(consts.ErrAssistantNotFound).Status(http.StatusUnauthorized).Send()
		return
	}

	if err != nil {
		c.Abort()
		response.Error(err).Status(http.StatusUnauthorized).Send()
		return
	}

	c.Set(consts.AuthAssistantShareMiddlewareKey, &assistantEntity.Assistant)

	c.Next()
}

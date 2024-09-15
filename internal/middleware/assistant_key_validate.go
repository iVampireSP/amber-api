package middleware

import (
	"net/http"
	"rag-new/internal/schema"
	"rag-new/internal/service/assistant"
	"rag-new/pkg/consts"
	"strings"

	"github.com/gin-gonic/gin"
)

type AssistantKeyValidateMiddleware struct {
	assistantService *assistant.Service
}

func NewAssistantKeyValidateMiddleware(assistantService *assistant.Service) *AssistantKeyValidateMiddleware {
	return &AssistantKeyValidateMiddleware{
		assistantService,
	}
}

func (a *AssistantKeyValidateMiddleware) AssistantKeytValidate(c *gin.Context) {
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

	key := authSplit[1]
	// if start with sk-
	if strings.HasPrefix(key, "sk-") {
		key = key[3:]
	}

	if key == "" {
		c.Abort()
		response.Error(consts.ErrAssistantKeyNotFound).Status(http.StatusUnauthorized).Send()
		return
	}

	assistantEntity, err := a.assistantService.GetByKey(c, key)
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

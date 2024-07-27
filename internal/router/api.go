package router

import (
	"github.com/gin-gonic/gin"
	v1 "rag-new/internal/api/v1"
)

type Api struct {
	User *v1.UserController
}

func NewApiRoute(
	User *v1.UserController,
) *Api {
	return &Api{
		User: User,
	}
}

func (a *Api) InitApiRouter(r *gin.RouterGroup) {
	r.GET("/ping", a.User.Test)
}

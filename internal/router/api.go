package router

import (
	"github.com/gin-gonic/gin"
	v1 "rag-new/internal/api/v1"
)

type Api struct {
	User *v1.UserController
	Tool *v1.ToolController
}

func NewApiRoute(
	User *v1.UserController,
	Tool *v1.ToolController,
) *Api {
	return &Api{
		User, Tool,
	}
}

func (a *Api) InitApiRouter(r *gin.RouterGroup) {
	r.GET("/ping", a.User.Test)

	r.GET("/tools", a.Tool.List)
	r.POST("/tools", a.Tool.CreateTool)
}

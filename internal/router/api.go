package router

import (
	"github.com/gin-gonic/gin"
	v1 "rag-new/internal/api/v1"
)

type Api struct {
	User      *v1.UserController
	Tool      *v1.ToolController
	Assistant *v1.AssistantController
}

func NewApiRoute(
	User *v1.UserController,
	Tool *v1.ToolController,
	Assistant *v1.AssistantController,
) *Api {
	return &Api{
		User, Tool, Assistant,
	}
}

func (a *Api) InitApiRouter(r *gin.RouterGroup) {
	r.GET("/ping", a.User.Test)

	r.GET("/tools", a.Tool.List)
	r.POST("/tools", a.Tool.CreateTool)
	r.GET("/tools/:id", a.Tool.GetTool)
	r.DELETE("/tools/:id", a.Tool.DeleteTool)

	r.GET("/assistants", a.Assistant.List)
	r.POST("/assistants", a.Assistant.CreateAssistant)
	//r.GET("/assistants/:id", a.Assistant.List)
	//r.DELETE("/assistants/:id", a.Assistant.DeleteAssistant)

	r.GET("/assistants/:id/tools", a.Assistant.ListTool)
	r.POST("/assistants/:id/tools/:tool_id", a.Assistant.BindTool)
}

package router

import (
	"github.com/gin-gonic/gin"
	v1 "rag-new/internal/api/v1"
)

type Api struct {
	User      *v1.UserController
	Tool      *v1.ToolController
	Assistant *v1.AssistantController
	Chat      *v1.ChatController
}

func NewApiRoute(
	User *v1.UserController,
	Tool *v1.ToolController,
	Assistant *v1.AssistantController,
	Chat *v1.ChatController,
) *Api {
	return &Api{
		User, Tool, Assistant, Chat,
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
	r.DELETE("/assistants/:id", a.Assistant.DeleteAssistant)

	r.GET("/assistants/:id/tools", a.Assistant.ListTool)
	r.POST("/assistants/:id/tools/:tool_id", a.Assistant.BindTool)
	r.DELETE("/assistants/:id/tools/:tool_id", a.Assistant.UnbindTool)

	r.GET("/chats", a.Chat.List)
	r.POST("/chats", a.Chat.Create)
	r.DELETE("/chats/:id", a.Chat.Delete)

	r.GET("/chats/:id/messages", a.Chat.ListChatMessage)
	r.POST("/chats/:id/messages", a.Chat.AddChatMessage)
}

func (a *Api) InitNoAuthApiRouter(r *gin.RouterGroup) {
	r.GET("/stream/:stream_id", a.Chat.Stream)
}

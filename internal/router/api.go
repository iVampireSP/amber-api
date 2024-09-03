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
	File      *v1.FileController
	Memory    *v1.MemoryController
}

func NewApiRoute(
	User *v1.UserController,
	Tool *v1.ToolController,
	Assistant *v1.AssistantController,
	Chat *v1.ChatController,
	File *v1.FileController,
	Memory *v1.MemoryController,
) *Api {
	return &Api{
		User, Tool, Assistant, Chat, File, Memory,
	}
}

func (a *Api) InitApiRouter(r *gin.RouterGroup) {
	r.GET("/ping", a.User.Test)

	r.GET("/assistants", a.Assistant.List)
	r.GET("/assistants/:id", a.Assistant.GetAssistant)
	r.PATCH("/assistants/:id", a.Assistant.UpdateAssistant)
	r.POST("/assistants", a.Assistant.CreateAssistant)
	r.DELETE("/assistants/:id", a.Assistant.DeleteAssistant)
	r.GET("/assistants/:id/shares", a.Assistant.ListAssistantShares)
	r.POST("/assistants/:id/shares", a.Assistant.CreateAssistantShare)
	r.DELETE("/assistants/:id/shares/:share_id", a.Assistant.DeleteAssistantShare)

	r.GET("/tools", a.Tool.List)
	r.POST("/tools", a.Tool.CreateTool)
	r.GET("/tools/:id", a.Tool.GetTool)
	r.DELETE("/tools/:id", a.Tool.DeleteTool)
	r.POST("/tools/:id/update", a.Tool.UpdateToolData)
	r.POST("/tools/syntax", a.Tool.ValidateSyntax)

	r.GET("/assistants/:id/tools", a.Assistant.ListTool)
	r.POST("/assistants/:id/tools/:tool_id", a.Assistant.BindTool)
	r.DELETE("/assistants/:id/tools/:tool_id", a.Assistant.UnbindTool)

	r.GET("/chats", a.Chat.List)
	r.POST("/chats", a.Chat.Create)
	r.POST("/chats/:id/clear", a.Chat.ClearChatMessage)
	r.DELETE("/chats/:id", a.Chat.Delete)

	r.GET("/chats/:id/messages", a.Chat.ListChatMessage)
	r.POST("/chats/:id/messages", a.Chat.AddChatMessage)
	r.POST("/chats/:id/images", a.Chat.AddChatImage)

	r.GET("/memories", a.Memory.List)
	r.POST("/memories/purge", a.Memory.Purge)
	r.DELETE("/memories/:id", a.Memory.Delete)
}

func (a *Api) InitNoAuthApiRouter(r *gin.RouterGroup) {
	r.GET("/stream/:stream_id", a.Chat.Stream)
	r.GET("/chat_public", a.Chat.GetChatPublic)
	r.POST("/chat_public", a.Chat.CreatePublicChat)
	r.GET("/chat_public/:chat_id/messages", a.Chat.GetPublicChatMessages)
	r.POST("/chat_public/:chat_id/messages", a.Chat.AddPublicChatMessages)
	r.POST("/chat_public/:chat_id/clear", a.Chat.ClearPublicChatMessages)
	r.POST("/chat_public/:chat_id/images", a.Chat.AddPublicChatImage)

	r.GET("/files/:id/download", a.File.DownloadFile)
}

func (a *Api) InitOpenAICompatibleApiRouter(r *gin.RouterGroup) {
	r.POST("/chat/completions", a.Chat.OpenAIChatCompletion)
}

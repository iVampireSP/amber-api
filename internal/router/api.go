package router

import (
	v1 "rag-new/internal/api/v1"

	"github.com/gin-gonic/gin"
)

type Api struct {
	User      *v1.UserController
	Tool      *v1.ToolController
	Assistant *v1.AssistantController
	Chat      *v1.ChatController
	File      *v1.FileController
	Memory    *v1.MemoryController
	Library   *v1.LibraryController
}

func NewApiRoute(
	User *v1.UserController,
	Tool *v1.ToolController,
	Assistant *v1.AssistantController,
	Chat *v1.ChatController,
	File *v1.FileController,
	Memory *v1.MemoryController,
	Library *v1.LibraryController,
) *Api {
	return &Api{
		User, Tool, Assistant, Chat, File, Memory, Library,
	}
}

func (a *Api) InitApiRouter(r *gin.RouterGroup) {
	r.GET("/ping", a.User.Test)

	r.GET("/assistants", a.Assistant.List)
	r.GET("/assistants/:id", a.Assistant.GetAssistant)
	r.PUT("/assistants/:id", a.Assistant.UpdateAssistant)
	r.POST("/assistants", a.Assistant.CreateAssistant)
	r.DELETE("/assistants/:id", a.Assistant.DeleteAssistant)
	r.GET("/assistants/:id/keys", a.Assistant.ListAssistantKeys)
	r.POST("/assistants/:id/keys", a.Assistant.CreateAssistantKey)
	r.DELETE("/assistants/:id/keys/:key_id", a.Assistant.DeleteAssistantKey)
	r.GET("/assistants/public", a.Assistant.AssistantPublicList)
	r.POST("/assistants/public/:id", a.Assistant.FavoriteAssistant)
	r.DELETE("/assistants/public/:id", a.Assistant.UnFavoriteAssistant)
	r.GET("/assistants/favorites", a.Assistant.FavoriteAssistants)

	r.GET("/tools", a.Tool.List)
	r.POST("/tools", a.Tool.CreateTool)
	r.GET("/tools/:id", a.Tool.GetTool)
	r.DELETE("/tools/:id", a.Tool.DeleteTool)
	r.POST("/tools/:id/update", a.Tool.UpdateToolData)
	r.POST("/tools/syntax", a.Tool.ValidateSyntax)

	r.GET("/assistants/:id/tools", a.Assistant.ListTool)
	r.POST("/assistants/:id/tools/:tool_id", a.Assistant.BindTool)
	r.DELETE("/assistants/:id/tools/:tool_id", a.Assistant.UnbindTool)

	r.POST("/assistants/:id/library", a.Assistant.BindLibrary)
	r.DELETE("/assistants/:id/library", a.Assistant.UnbindLibrary)

	r.GET("/chats", a.Chat.List)
	r.POST("/chats", a.Chat.Create)
	r.GET("/chats/:id", a.Chat.Show)
	r.PUT("/chats/:id", a.Chat.Update)
	r.POST("/chats/:id/clear", a.Chat.ClearChatMessage)
	r.DELETE("/chats/:id", a.Chat.Delete)

	r.GET("/chats/:id/messages", a.Chat.ListChatMessage)
	r.POST("/chats/:id/messages", a.Chat.AddChatMessage)
	r.POST("/chats/:id/files", a.Chat.AddChatFile)

	r.GET("/memories", a.Memory.List)
	r.POST("/memories/purge", a.Memory.Purge)
	r.DELETE("/memories/:id", a.Memory.Delete)

	r.GET("/libraries", a.Library.List)
	r.GET("/libraries/:id", a.Library.GetLibrary)
	r.POST("/libraries", a.Library.CreateLibrary)
	r.PUT("/libraries/:id", a.Library.Update)
	r.DELETE("/libraries/:id", a.Library.Delete)
	r.GET("/libraries/:id/documents", a.Library.ListDocuments)
	r.POST("/libraries/:id/documents", a.Library.CreateDocument)
	//r.PATCH("/libraries/:id/documents/:document_id", a.Library.UpdateDocument)
	r.DELETE("/libraries/:id/documents/:document_id", a.Library.DeleteDocument)
}

func (a *Api) InitNoAuthApiRouter(r *gin.RouterGroup) {
	r.GET("/stream/:stream_id", a.Chat.Stream)
	r.GET("/chat_public", a.Chat.GetChatPublic)
	r.POST("/chat_public", a.Chat.CreatePublicChat)
	r.GET("/chat_public/:chat_id/messages", a.Chat.GetPublicChatMessages)
	r.POST("/chat_public/:chat_id/messages", a.Chat.AddPublicChatMessages)
	r.POST("/chat_public/:chat_id/clear", a.Chat.ClearPublicChatMessages)
	r.POST("/chat_public/:chat_id/images", a.Chat.AddPublicChatImage)

	r.GET("/files/download/:hash", a.File.DownloadImage)
	//r.GET("/files/user/:id/download", a.File.DownloadUserFile)

}

func (a *Api) InitOpenAICompatibleApiRouter(r *gin.RouterGroup) {
	r.POST("/chat/completions", a.Chat.OpenAIChatCompletion)
}

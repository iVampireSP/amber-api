package schema

type AssistantIDRequest struct {
	ID EntityId `uri:"id" binding:"required"`
}

type AssistantCreateRequest struct {
	Name                 string  `json:"name" binding:"required" validate:"max=255"`
	Description          string  `json:"description" binding:"required" validate:"max=255"`
	Prompt               string  `json:"prompt" validate:"max=512"`
	DisableDefaultPrompt bool    `json:"disable_default_prompt" validate:"oneof=true false"`
	UserId               UserId  `json:"user_id" swaggerignore:"true" binding:"-"`
	Temperature          float64 `json:"temperature" validate:"oneof=0.1 0.2 0.3 0.4 0.5 0.6 0.7 0.8 0.9 1"`
}

type AssistantUpdateRequest struct {
	Name                        string    `json:"name" validate:"max=255"`
	Description                 string    `json:"description" validate:"max=255"`
	Prompt                      string    `json:"prompt" validate:"max=512"`
	DisableDefaultPrompt        bool      `json:"disable_default_prompt" validate:"oneof=true false"`
	DisableWebBrowsing          bool      `json:"disable_web_browsing" validate:"oneof=true false"`
	DisableMemory               bool      `json:"disable_memory" validate:"oneof=true false"`
	EnableMemoryForAssistantAPI bool      `json:"enable_memory_for_assistant_api" validate:"oneof=true false"`
	Public                      bool      `json:"public" validate:"oneof=true false"`
	LibraryId                   *EntityId `json:"library_id"`
	Temperature                 float64   `json:"temperature" validate:"oneof=0.1 0.2 0.3 0.4 0.5 0.6 0.7 0.8 0.9 1"`
}

type AssistantToolBindRequest struct {
	ToolId      EntityId `json:"tool_id" binding:"required"`
	AssistantId EntityId `json:"assistant_id" binding:"required"`
}

type AssistantToolUnbindRequest struct {
	ToolId      EntityId `json:"tool_id" binding:"required"`
	AssistantId EntityId `json:"assistant_id" binding:"required"`
}

type AssistantLibraryRequest struct {
	LibraryId EntityId `json:"library_id" binding:"required"`
}

type AssistantPublic struct {
	Id          EntityId `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
}

type PaginationRequest struct {
	Page int `form:"page"`
}

type AssistantSceneIdRequest struct {
	SceneId EntityId `uri:"scene_id" binding:"required"`
}

type CreateAssistantScenePromptRequest struct {
	Label  string `json:"label" binding:"required" validate:"max=20"`
	Prompt string `json:"prompt" binding:"required" validate:"max=512"`
}

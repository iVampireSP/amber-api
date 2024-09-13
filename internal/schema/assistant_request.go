package schema

type AssistantIDRequest struct {
	ID EntityId `uri:"id" binding:"required"`
}

type AssistantCreateRequest struct {
	Name                 string `json:"name" binding:"required" validate:"max=255"`
	Description          string `json:"description" binding:"required" validate:"max=255"`
	Prompt               string `json:"prompt" validate:"max=512"`
	DisableDefaultPrompt bool   `json:"disable_default_prompt" validate:"oneof=true false"`
	UserId               UserId `json:"user_id" swaggerignore:"true" binding:"-"`
}

type AssistantUpdateRequest struct {
	Name                          string    `json:"name" validate:"max=255"`
	Description                   string    `json:"description" validate:"max=255"`
	Prompt                        string    `json:"prompt" validate:"max=512"`
	DisableDefaultPrompt          bool      `json:"disable_default_prompt" validate:"oneof=true false"`
	DisableMemory                 bool      `json:"disable_memory" validate:"oneof=true false"`
	EnableMemoryForAssistantShare bool      `json:"enable_memory_for_assistant_share" validate:"oneof=true false"`
	LibraryId                     *EntityId `json:"library_id"`
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

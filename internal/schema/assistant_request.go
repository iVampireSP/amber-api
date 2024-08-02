package schema

type AssistantCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=255"`
	Description string `json:"description" binding:"required" validate:"max=255"`
	Prompt      string `json:"prompt"`
	UserId      UserId `json:"user_id" swaggerignore:"true" binding:"-"`
}

type AssistantUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

type AssistantToolBindRequest struct {
	ToolId      int64 `json:"tool_id" binding:"required"`
	AssistantId int64 `json:"assistant_id" binding:"required"`
}

type AssistantToolUnbindRequest struct {
	ToolId      int64 `json:"tool_id" binding:"required"`
	AssistantId int64 `json:"assistant_id" binding:"required"`
}

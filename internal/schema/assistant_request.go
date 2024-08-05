package schema

type AssistantCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=255"`
	Description string `json:"description" binding:"required" validate:"max=255"`
	Prompt      string `json:"prompt" validate:"max=512"`
	UserId      UserId `json:"user_id" swaggerignore:"true" binding:"-"`
}

type AssistantUpdateRequest struct {
	Name        string `json:"name" validate:"max=255"`
	Description string `json:"description" validate:"max=255"`
	Prompt      string `json:"prompt" validate:"max=512"`
}

type AssistantToolBindRequest struct {
	ToolId      int64 `json:"tool_id" binding:"required"`
	AssistantId int64 `json:"assistant_id" binding:"required"`
}

type AssistantToolUnbindRequest struct {
	ToolId      int64 `json:"tool_id" binding:"required"`
	AssistantId int64 `json:"assistant_id" binding:"required"`
}

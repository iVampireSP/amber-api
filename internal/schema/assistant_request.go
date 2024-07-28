package schema

type AssistantCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=255"`
	Description string `json:"description" binding:"required" validate:"max=255"`
	Prompt      string `json:"prompt"`
	UserId      UserId `json:"user_id" binding:"-"`
}

type AssistantUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

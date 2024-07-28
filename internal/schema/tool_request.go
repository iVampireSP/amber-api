package schema

type ToolCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=255"`
	Description string `json:"description" binding:"required" validate:"max=255"`
	Url         string `json:"url" binding:"required" validate:"max=255"`
	ApiKey      string `json:"api_key" validate:"max=255"`
}

type ToolUpdateRequest struct {
}

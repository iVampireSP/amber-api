package schema

type ChatCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=255"`
	AssistantId int64  `json:"assistant_id" binding:"required"`
	UserId      UserId `json:"user_id" binding:"-"`
}

type ChatHistoryAddRequest struct {
	Message string `json:"message" binding:"required" validate:"max=255"`
}

type ChatHistoryResponse struct {
	StreamId string `json:"stream_id"`
}

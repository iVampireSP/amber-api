package schema

type ChatCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=255"`
	AssistantId int64  `json:"assistant_id,string" binding:"required"`
	UserId      UserId `json:"user_id" swaggerignore:"true" binding:"-"`
}

type ChatGuestCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=255"`
	AssistantId int64  `json:"assistant_id,string" binding:"required"`
	GuestID     string `json:"guest_id" binding:"required" validate:"max=32"`
}

type ChatMessageAddRequest struct {
	Message string `json:"message" binding:"required" validate:"max=255"`
}

type ChatMessageResponse struct {
	StreamId string `json:"stream_id"`
}

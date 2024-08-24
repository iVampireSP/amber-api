package schema

import "mime/multipart"

type ChatRequest struct {
	ChatId int64 `uri:"id" binding:"required"`
}

type ChatCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=30"`
	AssistantId int64  `json:"assistant_id" binding:"required"`
	UserId      UserId `json:"user_id" swaggerignore:"true" binding:"-"`
}

type ChatGuestCreateRequest struct {
	Name        string `json:"name" binding:"required" validate:"max=30"`
	AssistantId int64  `json:"assistant_id" binding:"required"`
	GuestID     string `json:"guest_id" binding:"required" validate:"max=32"`
}

type ChatMessageAddRequest struct {
	Message string   `json:"message" binding:"required" validate:"max=255"`
	Role    ChatRole `json:"role" binding:"required" enums:"user,user_hide,system,system_hide,assistant,image"`
}

type ChatMessageAddImageRequest struct {
	Image *multipart.FileHeader `form:"image" binding:"required" swaggerignore:"true"`
}

type ChatMessageResponse struct {
	StreamId string `json:"stream_id"`
	Stream   bool   `json:"stream"`
}

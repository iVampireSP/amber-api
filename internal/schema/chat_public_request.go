package schema

type ChatPublicRequest struct {
	AssistantToken string `json:"assistant_token" binding:"required"`
	GuestId        string `json:"guest_id" binding:"required"`
	Name           string `json:"name" binding:"required" validate:"max=32"`
}

type ChatPublicResponse struct {
	AssistantToken string `json:"assistant_token" binding:"required"`
	GuestId        string `json:"guest_id" binding:"required"`
}

type ChatPublicListRequest struct {
	GuestId string `json:"guest_id" binding:"required"`
}

type GetPublicChatMessageRequestParams struct {
	ChatId int64 `uri:"chat_id" binding:"required"`
}

type GetPublicChatMessageRequest struct {
	AssistantToken string `query:"assistant_token" form:"assistant_token"  json:"assistant_token" binding:"required"`
	GuestId        string `query:"guest_id" form:"guest_id"  json:"guest_id" binding:"required"`
}

type AddPublicChatMessageRequest struct {
	AssistantToken string `json:"assistant_token" binding:"required"`
	GuestId        string `json:"guest_id" binding:"required"`
	Message        string `json:"message" binding:"required"`
}

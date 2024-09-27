package schema

type ChatPublicRequest struct {
	AssistantKey string `json:"assistant_key" binding:"required"`
	GuestId      string `json:"guest_id" binding:"required"`
	Name         string `json:"name" binding:"required" validate:"max=32"`
}

type ChatPublicResponse struct {
	AssistantKey string `json:"assistant_key" binding:"required"`
	GuestId      string `json:"guest_id" binding:"required"`
}

type ChatPublicListRequest struct {
	GuestId string `json:"guest_id" binding:"required"`
}

type GetPublicChatMessageRequestParams struct {
	ChatId EntityId `uri:"chat_id" binding:"required"`
}

type GetPublicChatMessageRequest struct {
	AssistantKey string `query:"assistant_key" form:"assistant_key"  json:"assistant_key" binding:"required"`
	GuestId      string `query:"guest_id" form:"guest_id"  json:"guest_id" binding:"required"`
}

type AddPublicChatMessageRequest struct {
	AssistantKey string   `json:"assistant_key" binding:"required"`
	GuestId      string   `json:"guest_id" binding:"required"`
	Message      string   `json:"message" binding:"required"`
	Role         ChatRole `json:"role" binding:"required" enums:"user,user_hide,system,system_hide,assistant,system_override,user_later"`
}

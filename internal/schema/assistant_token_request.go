package schema

type AssistantTokenListRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
}

type AssistantTokenCreateRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
}

type AssistantTokenDeleteRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
	TokenId     EntityId `uri:"share_id" binding:"required"`
}

type AssistantTokenUpdateRequest struct {
	AssistantId int64
	Token       string
}

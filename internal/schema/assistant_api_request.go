package schema

type AssistantKeyListRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
}

type AssistantKeyCreateRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
}

type AssistantKeyDeleteRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
	KeyId       EntityId `uri:"key_id" binding:"required"`
}

type AssistantKeyUpdateRequest struct {
	AssistantId int64
	Token       string
}

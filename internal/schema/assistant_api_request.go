package schema

type AssistantApiKeyListRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
}

type AssistantApiKeyCreateRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
}

type AssistantApiKeyDeleteRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
	KeyId       EntityId `uri:"key_id" binding:"required"`
}

type AssistantApiKeyUpdateRequest struct {
	AssistantId int64
	Token       string
}

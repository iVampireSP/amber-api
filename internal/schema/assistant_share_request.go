package schema

type AssistantShareListRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
}

type AssistantShareCreateRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
}

type AssistantShareDeleteRequest struct {
	AssistantId EntityId `uri:"id" binding:"required"`
	ShareId     EntityId `uri:"share_id" binding:"required"`
}

type AssistantShareUpdateRequest struct {
	AssistantId int64
	Token       string
}

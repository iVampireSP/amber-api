package schema

type AssistantShareListRequest struct {
	AssistantId int64 `uri:"id" binding:"required"`
}

type AssistantShareCreateRequest struct {
	AssistantId int64 `uri:"id" binding:"required"`
}

type AssistantShareDeleteRequest struct {
	AssistantId int64 `uri:"id" binding:"required"`
	ShareId     int64 `uri:"share_id" binding:"required"`
}

type AssistantShareUpdateRequest struct {
	AssistantId int64
	Token       string
}

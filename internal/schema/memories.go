package schema

type MemoryListRequest struct {
}

type MemoryDeleteRequest struct {
	ID EntityId `uri:"id" binding:"required"`
}

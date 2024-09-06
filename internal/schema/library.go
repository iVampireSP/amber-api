package schema

type LibraryIdRequest struct {
	Id EntityId `uri:"id"`
}

type LibraryUpdateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

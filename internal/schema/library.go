package schema

type LibraryIdRequest struct {
	Id EntityId `uri:"id"`
}

type LibraryAndDocumentIdRequest struct {
	Id         EntityId `uri:"id"`
	DocumentId EntityId `uri:"document_id"`
}

type LibraryUpdateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Default     *bool   `json:"default"`
}

type LibraryCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type DocumentUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

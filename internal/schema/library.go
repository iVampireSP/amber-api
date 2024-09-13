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

type DocumentCreateRequest struct {
	Name    string `json:"name" binding:"required" min:"1" max:"128"`
	Content string `json:"content" binding:"required" min:"1" max:"65536"`
}

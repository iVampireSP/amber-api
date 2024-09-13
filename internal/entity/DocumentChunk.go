package entity

import "rag-new/internal/schema"

type DocumentChunk struct {
	Model
	Content    string          `json:"content"`
	Order      int             `json:"order"`
	DocumentId schema.EntityId `json:"document_id"`
	Library    *Library        `json:"library"`
	LibraryId  schema.EntityId `json:"library_id"`
	// Vectorized 代表是否已经向量化，不应该用 Chunked
	Vectorized bool `json:"vectorized"`
}

func (*DocumentChunk) TableName() string {
	return "document_chunks"
}

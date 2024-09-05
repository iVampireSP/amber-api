package entity

import "rag-new/internal/schema"

type DocumentChunk struct {
	Model
	Content    string          `json:"content"`
	Order      int             `json:"order"`
	DocumentId schema.EntityId `json:"document_id"`
	Library    *Library        `json:"library"`
	LibraryId  schema.EntityId `json:"library_id"`
	Chunked    bool            `json:"chunked"`
}

func (*DocumentChunk) TableName() string {
	return "document_chunks"
}

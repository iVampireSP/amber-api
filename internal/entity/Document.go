package entity

import "rag-new/internal/schema"

type Document struct {
	Model
	Name      string          `json:"name"`
	Chunked   bool            `json:"chunked"`
	LibraryId schema.EntityId `json:"library_id"`
	Library   *Library        `json:"library"`
}

func (*Document) TableName() string {
	return "documents"
}

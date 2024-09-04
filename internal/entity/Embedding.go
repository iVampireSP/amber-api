package entity

import "rag-new/internal/schema"

type Embedding struct {
	Model

	Text           string           `json:"name"`
	FileId         *schema.EntityId `json:"file_id"`
	File           *File            `json:"file"`
	TextMd5        string           `json:"text_md5"`
	EmbeddingModel string           `gorm:"column:model" json:"model"`
	Vector         schema.Embedding `json:"vector"`
}

func (t *Embedding) TableName() string {
	return "embeddings"
}

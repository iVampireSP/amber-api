package entity

import (
	"rag-new/internal/schema"
)

type Memory struct {
	Model

	Content        string           `json:"content"`
	ContentMd5     string           `json:"-"`
	EmbeddingModel string           `gorm:"column:model" json:"-"`
	Vector         schema.Embedding `json:"-"`
	UserId         schema.UserId    `json:"user_id"`
}

func (t *Memory) TableName() string {
	return "memories"
}

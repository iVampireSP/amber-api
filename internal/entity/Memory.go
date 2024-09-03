package entity

import (
	"rag-new/internal/schema"
)

type Memory struct {
	Model

	Content    string           `json:"content"`
	ContentMd5 string           `json:"content_md5"`
	Vector     schema.Embedding `json:"vector"`
	UserId     schema.UserId    `json:"user_id"`
}

func (t *Memory) TableName() string {
	return "memories"
}

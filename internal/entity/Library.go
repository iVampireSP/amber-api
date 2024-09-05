package entity

import "rag-new/internal/schema"

type Library struct {
	Model
	Name        string        `json:"name"`
	Default     bool          `json:"default"`
	Description *string       `json:"description"`
	UserId      schema.UserId `json:"user_id"`
}

func (*Library) TableName() string {
	return "libraries"
}

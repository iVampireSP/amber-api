package entity

import "rag-new/internal/schema"

type UnsettledToken struct {
	Model
	UserId schema.UserId `json:"user_id"`
	Count  int64         `json:"count"`
}

func (*UnsettledToken) TableName() string {
	return "unsettled_tokens"
}

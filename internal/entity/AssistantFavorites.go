package entity

import "rag-new/internal/schema"

type AssistantFavorites struct {
	Model
	AssistantId schema.EntityId `json:"assistant_id"`
	Assistant   *Assistant      `json:"assistant"`
	UserId      schema.UserId   `json:"user_id"`
}

func (a *AssistantFavorites) TableName() string {
	return "assistant_favorites"
}

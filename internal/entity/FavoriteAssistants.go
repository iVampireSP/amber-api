package entity

import "rag-new/internal/schema"

type FavoriteAssistants struct {
	Model
	AssistantId schema.EntityId `json:"assistant_id"`
	Assistant   *Assistant      `json:"assistant"`
	UserId      schema.UserId   `json:"user_id"`
}

func (a *FavoriteAssistants) TableName() string {
	return "favorite_assistants"
}

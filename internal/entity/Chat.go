package entity

import (
	"rag-new/internal/schema"
	"time"
)

type Chat struct {
	Model
	Name        string           `json:"name"`
	AssistantId *schema.EntityId `json:"assistant_id"`
	Assistant   *Assistant       `json:"-"`
	Prompt      *string          `json:"prompt"`
	UserId      schema.UserId    `json:"user_id"`
	ExpiredAt   *time.Time       `json:"expired_at"`
	Owner       schema.ChatOwner `json:"owner"`
	GuestId     *string          `json:"guest_id"`
}

func (a *Model) TableName() string {
	return "chats"
}

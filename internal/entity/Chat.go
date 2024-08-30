package entity

import (
	"rag-new/internal/schema"
	"time"
)

type Chat struct {
	Model       `xorm:"extends"`
	Name        string          `xorm:"varchar(255) notnull" json:"name"`
	AssistantId schema.EntityId `xorm:"varchar(255) notnull" json:"assistant_id"`
	//Assistant   Assistant        `xorm:"extends" json:"assistant"`
	UserId    schema.UserId    `xorm:"user_id int(11)" json:"user_id"`
	ExpiredAt *time.Time       `xorm:"TIMESTAMP null" json:"expired_at"`
	Owner     schema.ChatOwner `xorm:"varchar(255) notnull" json:"owner"`
	GuestId   *string          `xorm:"varchar(255)" json:"guest_id"`
}

type ChatWithAssistant struct {
	Model
	Assistant *Assistant `xorm:"extends" json:"assistant"`
}

func (a *Model) TableName() string {
	return "chats"
}

type ChatMessage struct {
	Model            `xorm:"extends"`
	ChatId           schema.EntityId `xorm:"varchar(255) notnull" json:"assistant_id"`
	Content          string          `xorm:"varchar(255) notnull" json:"content"`
	Role             schema.ChatRole `xorm:"varchar(255) notnull" json:"role"`
	Hidden           bool            `xorm:"bool notnull" json:"hidden"`
	PromptTokens     int             `xorm:"INTEGER" json:"prompt_tokens"`
	CompletionTokens int             `xorm:"INTEGER" json:"completion_tokens"`
	TotalTokens      int             `xorm:"INTEGER" json:"total_tokens"`
}

func (at *ChatMessage) TableName() string {
	return "chat_messages"
}

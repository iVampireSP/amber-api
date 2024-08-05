package entity

import (
	"rag-new/internal/schema"
	"time"
)

type Chat struct {
	Base        `xorm:"extends"`
	Name        string           `xorm:"varchar(255) notnull" json:"name"`
	AssistantId int64            `xorm:"varchar(255) notnull" json:"assistant_id"`
	UserId      schema.UserId    `xorm:"user_id int(11)" json:"user_id"`
	ExpiredAt   time.Time        `xorm:"TIMESTAMP notnull" json:"expired_at"`
	Owner       schema.ChatOwner `xorm:"varchar(255) notnull" json:"owner"`
	GuestId     string           `xorm:"varchar(255)" json:"guest_id"`
}

type ChatWithAssistant struct {
	Chat      `xorm:"extends"`
	Assistant *Assistant `xorm:"extends" json:"assistant"`
}

func (a *Base) TableName() string {
	return "chats"
}

type ChatMessage struct {
	Base         `xorm:"extends"`
	ChatId       int64           `xorm:"varchar(255) notnull" json:"assistant_id"`
	Content      string          `xorm:"varchar(255) notnull" json:"content"`
	Role         schema.ChatRole `xorm:"varchar(255) notnull" json:"role"`
	InputTokens  int             `xorm:"INTEGER" json:"input_tokens"`
	OutputTokens int             `xorm:"INTEGER" json:"output_tokens"`
	TotalTokens  int             `xorm:"INTEGER" json:"total_tokens"`
}

func (at *ChatMessage) TableName() string {
	return "chat_messages"
}

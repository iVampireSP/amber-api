package entity

import (
	"rag-new/internal/schema"
)

type Chat struct {
	Base        `xorm:"extends"`
	Name        string        `xorm:"varchar(255) notnull" json:"name"`
	AssistantId int64         `xorm:"varchar(255) notnull" json:"assistant_id"`
	UserId      schema.UserId `xorm:"user_id int(11) notnull" json:"user_id"`
}

type ChatWithAssistant struct {
	Chat      `xorm:"extends"`
	Assistant *Assistant `xorm:"extends" json:"assistant"`
}

func (a *Base) TableName() string {
	return "chats"
}

type ChatRole string

var (
	RoleAssistant  ChatRole = "assistant"
	RoleHuman      ChatRole = "user"
	RoleSystem     ChatRole = "system"
	RoleHideSystem ChatRole = "system_hide"
)

func (c ChatRole) String() string {
	return string(c)
}

type ChatMessage struct {
	Base         `xorm:"extends"`
	ChatId       int64    `xorm:"varchar(255) notnull" json:"assistant_id"`
	Content      string   `xorm:"varchar(255) notnull" json:"content"`
	Role         ChatRole `xorm:"varchar(255) notnull" json:"role"`
	InputTokens  int      `xorm:"varchar(255) notnull" json:"input_tokens"`
	OutputTokens int      `xorm:"varchar(255) notnull" json:"output_tokens"`
	TotalTokens  int      `xorm:"varchar(255) notnull" json:"total_tokens"`
}

func (at *ChatMessage) TableName() string {
	return "chat_messages"
}

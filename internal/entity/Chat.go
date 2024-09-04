package entity

import (
	"rag-new/internal/schema"
	"time"
)

type Chat struct {
	Model
	Name        string           `json:"name"`
	AssistantId schema.EntityId  `json:"assistant_id"`
	Assistant   Assistant        `json:"-"`
	UserId      schema.UserId    ` json:"user_id"`
	ExpiredAt   *time.Time       `json:"expired_at"`
	Owner       schema.ChatOwner `json:"owner"`
	GuestId     *string          `json:"guest_id"`
}

type ChatWithAssistant struct {
	Model
	Assistant *Assistant ` json:"assistant"`
}

func (a *Model) TableName() string {
	return "chats"
}

type ChatMessage struct {
	Model
	ChatId           schema.EntityId  `json:"assistant_id"`
	Content          string           `json:"content"`
	Role             schema.ChatRole  `json:"role"`
	ToolCall         *schema.ToolCall `json:"-"`
	FileId           *schema.EntityId `json:"file_id"`
	File             *File            `json:"file"`
	Hidden           bool             `json:"hidden"`
	PromptTokens     int              `json:"prompt_tokens"`
	CompletionTokens int              `json:"completion_tokens"`
	TotalTokens      int              `json:"total_tokens"`
}

func (at *ChatMessage) TableName() string {
	return "chat_messages"
}

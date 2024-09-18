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
	UserId      schema.UserId    `json:"user_id"`
	ExpiredAt   *time.Time       `json:"expired_at"`
	Owner       schema.ChatOwner `json:"owner"`
	GuestId     *string          `json:"guest_id"`
}

func (a *Model) TableName() string {
	return "chats"
}

type ChatMessage struct {
	Model
	ChatId schema.EntityId `json:"chat_id"`
	// AssistantId 可以让同一个对话中，使用不同的助手来处理消息
	AssistantId *schema.EntityId `json:"assistant_id"`
	Assistant   *Assistant       `json:"assistant"`
	Content     string           `json:"content"`
	Role        schema.ChatRole  `json:"role"`
	ToolCall    *schema.ToolCall `json:"-"`
	// FileId
	FileId *schema.EntityId `json:"file_id"`
	File   *File            `json:"file"`
	//UserFileId       *schema.EntityId `json:"user_file_id"`
	//UserFile         *UserFile        `json:"user_file"`
	Hidden           bool `json:"hidden"`
	PromptTokens     int  `json:"prompt_tokens"`
	CompletionTokens int  `json:"completion_tokens"`
	TotalTokens      int  `json:"total_tokens"`
}

func (at *ChatMessage) TableName() string {
	return "chat_messages"
}

type ChatMessageList struct {
	Id          schema.EntityId  `json:"id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	ChatId      schema.EntityId  `json:"chat_id"`
	AssistantId *schema.EntityId `json:"assistant_id"`
	Assistant   *struct {
		Id   schema.EntityId `json:"id"`
		Name string          `json:"name"`
	} `json:"assistant"`
	Content string           `json:"content"`
	Role    schema.ChatRole  `json:"role"`
	FileId  *schema.EntityId `json:"file_id"`
	File    *File            `json:"file"`
	//UserFileId       *schema.EntityId `json:"user_file_id"`
	//UserFile         *UserFile        `json:"user_file"`
	Hidden           bool `json:"hidden"`
	PromptTokens     int  `json:"prompt_tokens"`
	CompletionTokens int  `json:"completion_tokens"`
	TotalTokens      int  `json:"total_tokens"`
}

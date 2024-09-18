package entity

import (
	"rag-new/internal/schema"
	"time"
)

type Tool struct {
	Model

	Name         string                     `json:"name"`
	Description  string                     `json:"description"`
	DiscoveryUrl string                     `json:"discovery_url"`
	ApiKey       string                     `json:"api_key"`
	Data         schema.ToolDiscoveryOutput `json:"data"`
	UserId       schema.UserId              `json:"user_id"`
}

func (t *Tool) TableName() string {
	return "tools"
}

type ToolCallToken struct {
	Model

	ChatId    schema.EntityId `json:"chat_id"`
	Token     string          `json:"token"`
	ExpiredAt time.Time       `json:"expired_at"`
}

func (t *ToolCallToken) TableName() string {
	return "tool_call_tokens"
}

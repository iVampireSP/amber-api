package entity

import "rag-new/internal/schema"

type AssistantApiKey struct {
	Model
	AssistantId schema.EntityId `json:"assistant_id"`
	Assistant   Assistant       `json:"assistant"`
	Secret      string          `json:"secret"`
}

func (a *AssistantApiKey) TableName() string {
	return "assistant_api_keys"
}

//func (a *AssistantShareType) TableName() string {
//	return "assistant_shares"
//}
//
//func (a *AssistantShareType) ToAssistantShare() *AssistantShare {
//	var assistantShare = &AssistantShare{}
//	assistantShare.Id = a.Id
//	assistantShare.AssistantId = a.AssistantId
//	assistantShare.Token = a.Token
//	assistantShare.CreatedAt = a.Assistant.CreatedAt
//	assistantShare.UpdatedAt = a.Assistant.UpdatedAt
//
//	return assistantShare
//
//}

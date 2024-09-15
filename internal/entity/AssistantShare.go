package entity

import "rag-new/internal/schema"

type AssistantShare struct {
	Model
	AssistantId schema.EntityId `json:"assistant_id"`
	Assistant   Assistant       `json:"assistant"`
	Token       string          `json:"token"`
}

func (a *AssistantShare) TableName() string {
	return "assistant_shares"
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

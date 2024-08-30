package entity

import "rag-new/internal/schema"

type AssistantShare struct {
	Model `xorm:"extends"`

	AssistantId schema.EntityId `xorm:"varchar(255) notnull" json:"assistant_id"`
	Token       string          `xorm:"varchar(255) notnull" json:"token"`
}

type AssistantShareType struct {
	Model `xorm:"extends"`

	Token       string          `xorm:"varchar(255) notnull" json:"token"`
	AssistantId schema.EntityId `xorm:"int(8) notnull index" json:"assistant_id"`
	Assistant   *Assistant      `xorm:"extends" json:"assistant"`
}

func (a *AssistantShare) TableName() string {
	return "assistant_shares"
}

func (a *AssistantShareType) TableName() string {
	return "assistant_shares"
}

func (a *AssistantShareType) ToAssistantShare() *AssistantShare {
	var assistantShare = &AssistantShare{}
	assistantShare.Id = a.Id
	assistantShare.AssistantId = a.AssistantId
	assistantShare.Token = a.Token
	assistantShare.CreatedAt = a.Assistant.CreatedAt
	assistantShare.UpdatedAt = a.Assistant.UpdatedAt

	return assistantShare

}

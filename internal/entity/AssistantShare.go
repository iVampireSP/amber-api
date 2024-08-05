package entity

type AssistantShare struct {
	Base        `xorm:"extends"`
	AssistantId int64  `xorm:"varchar(255) notnull" json:"assistant_id"`
	Token       string `xorm:"varchar(255) notnull" json:"token"`
}

type AssistantShareType struct {
	Base        `xorm:"extends"`
	Token       string     `xorm:"varchar(255) notnull" json:"token"`
	AssistantId int64      `xorm:"int(8) notnull index" json:"assistant_id"`
	Assistant   *Assistant `xorm:"extends" json:"assistant"`
}

func (a *AssistantShare) TableName() string {
	return "assistant_shares"
}

func (a *AssistantShareType) TableName() string {
	return "assistant_shares"
}

func (a *AssistantShareType) ToAssistantShare() *AssistantShare {
	var assistantShare = &AssistantShare{}
	assistantShare.ID = a.ID
	assistantShare.AssistantId = a.AssistantId
	assistantShare.Token = a.Token
	assistantShare.CreatedAt = a.Assistant.CreatedAt
	assistantShare.UpdatedAt = a.Assistant.UpdatedAt

	return assistantShare

}

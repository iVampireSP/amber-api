package entity

import (
	"rag-new/internal/schema"
)

type Assistant struct {
	Base        `xorm:"extends"`
	Name        string        `xorm:"varchar(255) notnull" json:"name"`
	Prompt      string        `xorm:"varchar(255) notnull" json:"prompt"`
	Description string        `xorm:"varchar(255) notnull" json:"description"`
	UserId      schema.UserId `xorm:"user_id int(11) notnull" json:"user_id"`
}

func (a *Assistant) TableName() string {
	return "assistants"
}

type AssistantTool struct {
	Base        `xorm:"extends"`
	AssistantId int64 `xorm:"varchar(255) notnull" json:"assistant_id"`
	ToolId      int64 `xorm:"varchar(255) notnull" json:"tool_id"`
}

func (at *AssistantTool) TableName() string {
	return "assistants"
}

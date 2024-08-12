package entity

import (
	"rag-new/internal/schema"
)

type Assistant struct {
	Base                 `xorm:"extends"`
	Name                 string        `xorm:"varchar(255) notnull" json:"name"`
	Prompt               string        `xorm:"varchar(255) notnull" json:"prompt"`
	Description          string        `xorm:"varchar(255) notnull" json:"description"`
	UserId               schema.UserId `xorm:"user_id int(11) notnull" json:"user_id"`
	DisableDefaultPrompt bool          `xorm:"disable_default_prompt bool notnull" json:"disable_default_prompt"`
}

func (a *Assistant) TableName() string {
	return "assistants"
}
func (a *Assistant) GetUserId() schema.UserId {
	return a.UserId
}

type AssistantTool struct {
	Base        `xorm:"extends"`
	AssistantId int64 `xorm:"int(255) notnull" json:"assistant_id"`
	ToolId      int64 `xorm:"int(255) notnull" json:"tool_id"`
}

type AssistantToolType struct {
	Base        `xorm:"extends"`
	AssistantId int64      `xorm:"int(8) notnull index" json:"assistant_id"`
	ToolId      int64      `xorm:"int(8) notnull index" json:"tool_id"`
	Assistant   *Assistant `xorm:"extends" json:"assistant"`
	Tool        *Tool      `xorm:"extends" json:"tool"`
}

func (at *AssistantTool) TableName() string {
	return "assistant_tools"
}
func (att *AssistantToolType) TableName() string {
	return "assistant_tools"
}

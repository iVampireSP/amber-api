package entity

import "rag-new/internal/schema"

type AssistantTool struct {
	Model
	AssistantId schema.EntityId `json:"assistant_id"`
	ToolId      schema.EntityId `json:"tool_id"`
	Assistant   Assistant       `json:"assistant"`
	Tool        Tool            `json:"tool"`
}

func (at *AssistantTool) TableName() string {
	return "assistant_tools"
}

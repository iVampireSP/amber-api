package entity

import (
	"rag-new/internal/schema"
)

type Assistant struct {
	Model
	Name                          string        `json:"name"`
	Prompt                        string        `json:"prompt"`
	Description                   string        `json:"description"`
	UserId                        schema.UserId `json:"user_id"`
	DisableDefaultPrompt          bool          `json:"disable_default_prompt"`
	DisableMemory                 bool          `json:"disable_memory"`
	EnableMemoryForAssistantShare bool          `json:"enable_memory_for_assistant_share"`
}

func (a *Assistant) TableName() string {
	return "assistants"
}
func (a *Assistant) GetUserId() schema.UserId {
	return a.UserId
}

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

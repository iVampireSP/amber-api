package entity

import "rag-new/internal/schema"

type ScenePrompt struct {
	Model
	Label       string          `json:"label"`
	Prompt      string          `json:"prompt"`
	AssistantId schema.EntityId `json:"assistant_id"`
	Assistant   *Assistant      `json:"assistant"`
}

func (*ScenePrompt) TableName() string {
	return "scene_prompts"
}

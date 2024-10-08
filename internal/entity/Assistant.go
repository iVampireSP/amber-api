package entity

import (
	"rag-new/internal/schema"
)

type Assistant struct {
	Model
	Name                        string           `json:"name"`
	Prompt                      string           `json:"prompt"`
	Description                 string           `json:"description"`
	UserId                      schema.UserId    `json:"user_id"`
	TotalTokenUsage             int64            `json:"total_token_usage"`
	LibraryId                   *schema.EntityId `json:"library_id"`
	Library                     *Library         `json:"-"`
	Temperature                 float64          `json:"temperature"`
	Public                      bool             `json:"public"`
	DisableDefaultPrompt        bool             `json:"disable_default_prompt"`
	DisableMemory               bool             `json:"disable_memory"`
	EnableMemoryForAssistantAPI bool             `json:"enable_memory_for_assistant_api"`
}

func (a *Assistant) TableName() string {
	return "assistants"
}
func (a *Assistant) GetUserId() schema.UserId {
	return a.UserId
}

// ToPublic 转换为公开助理
func (a *Assistant) ToPublic() *schema.AssistantPublic {
	return &schema.AssistantPublic{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description,
	}
}

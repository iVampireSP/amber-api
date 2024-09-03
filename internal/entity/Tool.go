package entity

import (
	"rag-new/internal/schema"
)

type Tool struct {
	Model

	Name         string                     `json:"name"`
	Description  string                     `json:"description"`
	DiscoveryUrl string                     `json:"discovery_url"`
	ApiKey       string                     `json:"api_key"`
	Data         schema.ToolDiscoveryOutput `json:"data"`
	UserId       schema.UserId              `json:"user_id"`
}

func (t *Tool) TableName() string {
	return "tools"
}

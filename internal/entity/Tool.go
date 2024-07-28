package entity

import (
	"encoding/json"
	"rag-new/internal/schema"
)

type Tool struct {
	Base         `xorm:"extends"`
	Name         string        `xorm:"varchar(255) notnull" json:"name"`
	Description  string        `xorm:"varchar(255) notnull" json:"description"`
	DiscoveryUrl string        `xorm:"varchar(255) notnull" json:"discovery_url"`
	ApiKey       string        `xorm:"varchar(255) notnull" json:"api_key"`
	Data         ToolData      `xorm:"json" json:"data"`
	UserID       schema.UserId `xorm:"user_id int(11) notnull" json:"user_id"`
}

func (t *Tool) TableName() string {
	return "tools"
}

type ToolData struct {
	Name          string `json:"name"`
	HomepageUrl   string `json:"homepage_url"`
	CallbackUrl   string `json:"callback_url"`
	Description   string `json:"description"`
	ToolId        int    `json:"tool_id"`
	ToolFunctions []struct {
		Type     string `json:"type"`
		Function struct {
			Name        string      `json:"name"`
			Description string      `json:"description"`
			Parameters  interface{} `json:"parameters"`
			Required    []string    `json:"required"`
		} `json:"function"`
	} `json:"tool_functions"`
}

func (td *ToolData) FromDB(data []byte) error {
	return json.Unmarshal(data, td)
}

func (td *ToolData) ToDB() ([]byte, error) {
	return json.Marshal(td)
}

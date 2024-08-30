package entity

import (
	"rag-new/internal/schema"
)

type Tool struct {
	Model `xorm:"extends"`

	Name         string                     `xorm:"varchar(255) notnull" json:"name"`
	Description  string                     `xorm:"varchar(255) notnull" json:"description"`
	DiscoveryUrl string                     `xorm:"varchar(255) notnull" json:"discovery_url"`
	ApiKey       string                     `xorm:"varchar(255) notnull" json:"api_key"`
	Data         schema.ToolDiscoveryOutput `xorm:"json" json:"data"`
	UserId       schema.UserId              `xorm:"user_id int(11) notnull" json:"user_id"`
}

func (t *Tool) TableName() string {
	return "tools"
}

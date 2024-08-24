package entity

import "rag-new/internal/schema"

type Base struct {
	Id        int64  `json:"id"`
	CreatedAt string `xorm:"created timestamp" json:"created_at"`
	UpdatedAt string `xorm:"updated timestamp" json:"updated_at"`
}

// Entity 是所有 entity 的基类，后期要将所有的 Base 改成这种形式
type Entity struct {
	Id        schema.EntityId `json:"string,id"`
	CreatedAt string          `xorm:"created timestamp" json:"created_at"`
	UpdatedAt string          `xorm:"updated timestamp" json:"updated_at"`
}

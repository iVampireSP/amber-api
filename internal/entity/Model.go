package entity

import (
	"gorm.io/gorm"
	"rag-new/internal/schema"
	"time"
)

type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Model 是所有 entity 的基类，后期要将所有的 Base 改成这种形式
type Model struct {
	Id        schema.EntityId `gorm:"primarykey"  json:"string,id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index"`
}

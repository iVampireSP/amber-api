package entity

import (
	"rag-new/internal/schema"
	"time"
)

//type Base struct {
//	ID        int64          `gorm:"primarykey" json:"id"`
//	CreatedAt time.Time      `json:"created_at"`
//	UpdatedAt time.Time      `json:"updated_at"`
//	DeletedAt gorm.DeletedAt `gorm:"index"`
//}

// Model 是所有 entity 的基类，后期要将所有的 Base 改成这种形式
type Model struct {
	Id        schema.EntityId `gorm:"primarykey" json:"id,string"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	//DeletedAt gorm.DeletedAt  `gorm:"index"`
}

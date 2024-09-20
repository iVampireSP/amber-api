package entity

import "rag-new/internal/schema"

type MessageBlock struct {
	Model
	FullContent string
	ChatId      schema.EntityId
	Hash        string
	Message     []*ChatMessage `gorm:"column:messages;serializer:json"`
	// Temp 指当前消息是否是临时的，不完整的 block。不会保存到数据库，只在处理时需要
	Temp bool `gorm:"-"`
}

func (mb *MessageBlock) TableName() string {
	return "message_blocks"
}

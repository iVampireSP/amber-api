package chat

import (
	"gorm.io/gorm"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/chat_message"
	"xorm.io/xorm"
)

type Service struct {
	x  *xorm.Engine
	db *gorm.DB
	a  *assistant.Service
	cm *chat_message.Service
}

func NewService(x *xorm.Engine, db *gorm.DB, a *assistant.Service, cm *chat_message.Service) *Service {
	return &Service{x, db, a, cm}
}

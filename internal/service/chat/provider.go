package chat

import (
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/chat_message"
	"xorm.io/xorm"
)

type Service struct {
	x  *xorm.Engine
	a  *assistant.Service
	cm *chat_message.Service
}

func NewService(x *xorm.Engine, a *assistant.Service, cm *chat_message.Service) *Service {
	return &Service{x, a, cm}
}

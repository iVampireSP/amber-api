package chat

import (
	"rag-new/internal/dao"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/chat_message"
)

type Service struct {
	dao *dao.Query
	a   *assistant.Service
	cm  *chat_message.Service
}

func NewService(dao *dao.Query, a *assistant.Service, cm *chat_message.Service) *Service {
	return &Service{dao, a, cm}
}

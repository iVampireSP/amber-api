package chat_message

import (
	"rag-new/internal/dao"
	"rag-new/internal/service/message_block"
)

type Service struct {
	dao          *dao.Query
	messageBlock *message_block.Service
}

func NewService(
	dao *dao.Query,
	messageBlock *message_block.Service,
) *Service {
	return &Service{dao, messageBlock}
}

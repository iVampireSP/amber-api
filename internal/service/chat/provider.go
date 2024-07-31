package chat

import (
	"rag-new/internal/service/assistant"
	"xorm.io/xorm"
)

type Service struct {
	x *xorm.Engine
	a *assistant.Service
}

func NewService(x *xorm.Engine, a *assistant.Service) *Service {
	return &Service{x, a}
}

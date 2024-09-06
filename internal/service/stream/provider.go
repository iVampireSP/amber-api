package stream

import (
	"rag-new/internal/base/conf"
)

type Service struct {
	config *conf.Config
}

func NewService(config *conf.Config) *Service {
	return &Service{
		config,
	}
}

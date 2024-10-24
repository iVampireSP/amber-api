package text_classification

import (
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
)

type Service struct {
	config *conf.Config
	logger *logger.Logger
}

func NewService(config *conf.Config, logger *logger.Logger) *Service {
	return &Service{
		config,
		logger,
	}
}

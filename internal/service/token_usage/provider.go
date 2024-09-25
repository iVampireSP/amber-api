package token_usage

import (
	"rag-new/internal/base/logger"
	"rag-new/internal/base/redis"
)

type Service struct {
	redis  *redis.Redis
	logger *logger.Logger
}

func NewService(client *redis.Redis, logger *logger.Logger) *Service {
	return &Service{
		redis:  client,
		logger: logger,
	}
}

package file

import (
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/base/redis"
	"rag-new/internal/base/s3"
	"rag-new/internal/dao"
)

type Service struct {
	s3     *s3.S3
	config *conf.Config
	dao    *dao.Query
	redis  *redis.Redis
	logger *logger.Logger
}

func NewService(s3 *s3.S3, config *conf.Config, dao *dao.Query, redis *redis.Redis, logger *logger.Logger) *Service {
	return &Service{s3, config, dao, redis, logger}
}

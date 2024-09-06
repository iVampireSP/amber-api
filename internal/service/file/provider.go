package file

import (
	"github.com/redis/go-redis/v9"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/s3"
	"rag-new/internal/dao"
)

type Service struct {
	s3     *s3.S3
	config *conf.Config
	dao    *dao.Query
	redis  *redis.Client
}

func NewService(s3 *s3.S3, config *conf.Config, dao *dao.Query, redis *redis.Client) *Service {
	return &Service{s3, config, dao, redis}
}

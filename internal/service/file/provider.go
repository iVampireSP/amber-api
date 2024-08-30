package file

import (
	"rag-new/internal/base/conf"
	"rag-new/internal/base/s3"
	"rag-new/internal/dao"
)

type Service struct {
	s3     *s3.S3
	config *conf.Config
	dao    *dao.Query
}

func NewService(s3 *s3.S3, config *conf.Config, dao *dao.Query) *Service {
	return &Service{s3, config, dao}
}

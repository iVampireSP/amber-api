package file

import (
	"rag-new/internal/base/conf"
	"rag-new/internal/base/s3"
	"xorm.io/xorm"
)

type Service struct {
	s3     *s3.S3
	x      *xorm.Engine
	config *conf.Config
}

func NewService(s3 *s3.S3, x *xorm.Engine, config *conf.Config) *Service {
	return &Service{s3, x, config}
}

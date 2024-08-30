package file

import (
	"gorm.io/gorm"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/s3"
	"xorm.io/xorm"
)

type Service struct {
	s3     *s3.S3
	x      *xorm.Engine
	config *conf.Config
	db     *gorm.DB
}

func NewService(s3 *s3.S3, x *xorm.Engine, config *conf.Config, db *gorm.DB) *Service {
	return &Service{s3, x, config, db}
}

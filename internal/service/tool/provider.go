package tool

import (
	"rag-new/internal/base/conf"
	"xorm.io/xorm"
)

type Service struct {
	x      *xorm.Engine
	config *conf.Config
}

func NewService(x *xorm.Engine, config *conf.Config) *Service {
	return &Service{
		x,
		config,
	}
}

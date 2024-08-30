package tool

import (
	"gorm.io/gorm"
	"rag-new/internal/base/conf"
	"xorm.io/xorm"
)

type Service struct {
	x      *xorm.Engine
	db     *gorm.DB
	config *conf.Config
}

func NewService(x *xorm.Engine, config *conf.Config, db *gorm.DB) *Service {
	return &Service{
		x,
		db,
		config,
	}
}

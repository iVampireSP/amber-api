package tool

import (
	"rag-new/internal/base/conf"
	"rag-new/internal/dao"
)

type Service struct {
	dao    *dao.Query
	config *conf.Config
}

func NewService(config *conf.Config, dao *dao.Query) *Service {
	return &Service{
		dao,
		config,
	}
}

package unsettled_token

import (
	"rag-new/internal/base/conf"
	"rag-new/internal/dao"
)

type Service struct {
	dao    *dao.Query
	config *conf.Config
}

func NewService(
	dao *dao.Query,
	config *conf.Config,
) *Service {
	return &Service{dao, config}
}

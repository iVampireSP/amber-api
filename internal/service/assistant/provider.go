package assistant

import (
	"rag-new/internal/dao"
)

type Service struct {
	dao *dao.Query
}

func NewService(
	dao *dao.Query,
) *Service {
	return &Service{dao}
}

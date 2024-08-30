package assistant

import (
	"gorm.io/gorm"
	"rag-new/internal/dao"
	"xorm.io/xorm"
)

type Service struct {
	x   *xorm.Engine
	db  *gorm.DB
	dao *dao.Query
}

func NewService(
	x *xorm.Engine,
	db *gorm.DB,
	dao *dao.Query,
) *Service {
	return &Service{x, db, dao}
}

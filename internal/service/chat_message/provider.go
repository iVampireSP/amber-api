package chat_message

import (
	"gorm.io/gorm"
	"xorm.io/xorm"
)

type Service struct {
	x  *xorm.Engine
	db *gorm.DB
}

func NewService(x *xorm.Engine, db *gorm.DB) *Service {
	return &Service{x, db}
}

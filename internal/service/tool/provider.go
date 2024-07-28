package tool

import "xorm.io/xorm"

type Service struct {
	x *xorm.Engine
}

func NewService(x *xorm.Engine) *Service {
	return &Service{
		x: x,
	}
}

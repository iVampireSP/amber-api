package batch

import (
	"rag-new/internal/base/logger"
	"xorm.io/xorm"
)

type Batch struct {
	x      *xorm.Engine
	logger *logger.Logger
}

func NewBatch(
	x *xorm.Engine,
	logger *logger.Logger,
) *Batch {
	//base.NewApplication()
	return &Batch{x, logger}
}

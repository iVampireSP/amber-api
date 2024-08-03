package assistant

import (
	"rag-new/internal/batch"
	"xorm.io/xorm"
)

type Service struct {
	x     *xorm.Engine
	batch *batch.Batch
}

func NewService(
	x *xorm.Engine,
	batch *batch.Batch,
) *Service {
	return &Service{x, batch}
}

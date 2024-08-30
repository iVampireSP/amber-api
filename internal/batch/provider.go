package batch

import (
	"rag-new/internal/base/logger"
)

type Batch struct {
	logger *logger.Logger
}

func NewBatch(
	logger *logger.Logger,
) *Batch {
	//base.NewApplication()
	return &Batch{logger}
}

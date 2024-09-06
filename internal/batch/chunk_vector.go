package batch

import (
	"context"
	"rag-new/internal/dao"
	"rag-new/internal/service/library"
)

type ChunkVectorBatch struct {
	LibraryService *library.Service
	DAO            *dao.Query
}

func (*Batch) ChunkVector(ctx context.Context, cv *ChunkVectorBatch) error {
	chunks, i, err := cv.DAO.WithContext(ctx).DocumentChunk.Where(cv.DAO.DocumentChunk.Chunked.Is(false)).FindByPage(0, 10)
	if err != nil {
		return err
	}

	if i > 0 {
		for _, chunk := range chunks {
			err = cv.LibraryService.ChunkToMilvus(ctx, chunk)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

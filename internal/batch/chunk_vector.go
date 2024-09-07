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

func (b *Batch) ChunkVector(ctx context.Context, cv *ChunkVectorBatch) error {
	chunks, err := cv.DAO.WithContext(ctx).DocumentChunk.Where(cv.DAO.DocumentChunk.Chunked.Is(false)).Find()
	if err != nil {
		return err
	}

	if len(chunks) > 0 {
		for _, chunk := range chunks {
			b.logger.Sugar.Infof("Vectoring chunk: %v", chunk.Id)
			err = cv.LibraryService.ChunkToMilvus(ctx, chunk)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

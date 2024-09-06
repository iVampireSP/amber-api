package library

import (
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	entity2 "github.com/milvus-io/milvus-sdk-go/v2/entity"
	"rag-new/internal/entity"
)

func (s *Service) ChunkToMilvus(ctx context.Context, chunk *entity.DocumentChunk) error {
	emb, err := s.embedding.TextEmbedding(ctx, []string{chunk.Content})
	if err != nil {
		return err
	}

	var entityCols = []entity2.Column{
		entity2.NewColumnFloatVector("vector", s.config.OpenAI.EmbeddingDim, emb),
		entity2.NewColumnInt64("chunk_id", []int64{int64(chunk.Id)}),
		entity2.NewColumnInt64("library_id", []int64{int64(chunk.LibraryId)}),
		entity2.NewColumnInt64("document_id", []int64{int64(chunk.DocumentId)}),
		entity2.NewColumnVarChar("model", []string{s.config.OpenAI.EmbeddingModel}),
	}

	// insert to milvus
	_, err = s.milvus.Upsert(ctx, s.config.Milvus.DocumentCollection, "", entityCols...)
	if err != nil {
		return err
	}

	chunk.Chunked = true

	_, err = s.dao.WithContext(ctx).DocumentChunk.Where(s.dao.DocumentChunk.Id.Eq(uint(chunk.Id))).Update(s.dao.DocumentChunk.Chunked, true)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) SearchLibrary(ctx context.Context, content string, library *entity.Library) ([]*entity.DocumentChunk, error) {
	emb, err := s.embedding.TextEmbedding(ctx, []string{content})
	if err != nil {
		return nil, err
	}
	var filter = fmt.Sprintf("library_id == %d && model == %s", library.Id, s.config.OpenAI.EmbeddingModel)
	sp, err := entity2.NewIndexAUTOINDEXSearchParam(1)
	if err != nil {
		return nil, err
	}
	vector := entity2.FloatVector(emb[0])
	existingChunks, err := s.milvus.Search(ctx, s.config.Milvus.MemoryCollection,
		[]string{},
		filter,
		[]string{"document_id", "chunk_id"},
		[]entity2.Vector{vector},
		"vector",
		entity2.L2,
		10,
		sp, client.WithLimit(7))

	var ids []uint

	// get all data
	for _, res := range existingChunks {
		var chunkColumn *entity2.ColumnInt64
		for _, field := range res.Fields {
			if field.Name() == "chunk_id" {
				c, ok := field.(*entity2.ColumnInt64)
				if ok {
					chunkColumn = c
				}
			}
		}

		if chunkColumn == nil {
			return nil, fmt.Errorf("chunk column not found")
		}

		for i := 0; i < res.ResultCount; i++ {
			id, err := chunkColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}

			ids = append(ids, uint(id))

		}
	}

	documentChunks, err := s.dao.DocumentChunk.Where(s.dao.DocumentChunk.Where(s.dao.DocumentChunk.Id.In(ids...))).Find()

	return documentChunks, err
}

func (s *Service) deleteMilvusChunk(ctx context.Context, document *entity.Document) error {
	var filter = fmt.Sprintf("document_id == %d && model == %s", document.Id, s.config.OpenAI.EmbeddingModel)
	errDelete := s.milvus.Delete(ctx, s.config.Milvus.MemoryCollection, "", filter)
	return errDelete
}

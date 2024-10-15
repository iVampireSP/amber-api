package memory

import (
	"context"
	"fmt"
	"github.com/iVampireSP/pkg/md5"
	"rag-new/internal/entity"
	"rag-new/internal/schema"

	entity2 "github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func (s *Service) addMemory(ctx context.Context, data string, userId schema.UserId) (*entity.Memory, error) {
	dataMd5, err := md5.Md5(data)
	if err != nil {
		return nil, err
	}

	// check if memory exists
	count, err := s.dao.Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Where(s.dao.Memory.ContentMd5.Eq(dataMd5)).Where(s.dao.Memory.UserId.Eq(userId.String())).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		mem, err := s.dao.Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
			Where(s.dao.Memory.ContentMd5.Eq(dataMd5)).Where(s.dao.Memory.UserId.Eq(userId.String())).First()
		if err != nil {
			return nil, err
		}

		return mem, err
	}

	emb, err := s.Embedding.TextEmbedding(ctx, []string{data})
	if err != nil {
		return nil, err
	}

	var mem = &entity.Memory{
		Content:        data,
		ContentMd5:     dataMd5,
		Vector:         emb[0],
		UserId:         userId,
		EmbeddingModel: s.config.OpenAI.EmbeddingModel,
	}

	err = s.dao.Memory.Create(mem)
	if err != nil {
		return nil, err
	}

	var entityCols = []entity2.Column{
		entity2.NewColumnFloatVector("vector", s.config.OpenAI.EmbeddingDim, emb),
		entity2.NewColumnVarChar("hash", []string{dataMd5}),
		entity2.NewColumnVarChar("model", []string{s.config.OpenAI.EmbeddingModel}),
		entity2.NewColumnInt64("memory_id", []int64{int64(mem.Id)}),
		entity2.NewColumnVarChar("user_id", []string{userId.String()}),
	}

	// insert to milvus
	_, err = s.Milvus.Upsert(ctx, s.config.Milvus.MemoryCollection, "", entityCols...)
	if err != nil {
		return nil, err
	}

	return mem, err
}

func (s *Service) updateMemory(ctx context.Context, memoryId schema.EntityId, data string) (*entity.Memory, error) {
	dataMd5, err := md5.Md5(data)
	if err != nil {
		return nil, err
	}

	// 检查 memId 是否存在
	count, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Where(s.dao.Memory.Id.Eq(uint(memoryId))).Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		// 不存在
		return nil, fmt.Errorf("memory id not exists")
	}

	mem, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Where(s.dao.Memory.Id.Eq(uint(memoryId))).First()
	if err != nil {
		return nil, err
	}

	// 获取新的 embedding
	embed, err := s.Embedding.TextEmbedding(ctx, []string{data})
	if err != nil {
		return nil, err
	}

	mem.Vector = embed[0]
	// update
	_, err = s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Where(s.dao.Memory.Id.Eq(uint(memoryId))).Updates(mem)

	if err != nil {
		return nil, err
	}

	// 更新 milvus
	var entityCols = []entity2.Column{
		entity2.NewColumnFloatVector("vector", s.config.OpenAI.EmbeddingDim, embed),
		entity2.NewColumnVarChar("hash", []string{dataMd5}),
		entity2.NewColumnVarChar("model", []string{s.config.OpenAI.EmbeddingModel}),
		entity2.NewColumnInt64("memory_id", []int64{int64(memoryId)}),
		entity2.NewColumnVarChar("user_id", []string{mem.UserId.String()}),
	}

	// insert to milvus
	_, err = s.Milvus.Upsert(ctx, s.config.Milvus.MemoryCollection, "", entityCols...)
	if err != nil {
		return nil, err
	}

	return mem, nil

}

func (s *Service) deleteMemory(ctx context.Context, memoryId schema.EntityId) error {
	// 检查 memId 是否存在
	count, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Where(s.dao.Memory.Id.Eq(uint(memoryId))).Count()
	if err != nil {
		return err
	}

	if count == 0 {
		// 不存在
		return fmt.Errorf("delete failed, memory id not exists")
	}

	mem, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Where(s.dao.Memory.Id.Eq(uint(memoryId))).First()
	if err != nil {
		return err
	}

	// milvus delete
	ids := entity2.NewColumnInt64("memory_id", []int64{int64(mem.Id)})
	errDelete := s.Milvus.DeleteByPks(ctx, s.config.Milvus.MemoryCollection, "", ids)
	if errDelete != nil {
		return errDelete
	}

	// database delete
	_, err = s.dao.WithContext(ctx).Memory.Delete(mem)
	if err != nil {
		return err
	}

	return nil
}

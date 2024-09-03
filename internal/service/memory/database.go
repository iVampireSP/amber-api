package memory

import (
	"context"
	"fmt"
	entity2 "github.com/milvus-io/milvus-sdk-go/v2/entity"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/md5"
)

func (s *Service) addMemory(ctx context.Context, data string, userId schema.UserId) (schema.EntityId, error) {
	dataMd5, err := md5.Md5(data)
	if err != nil {
		return 0, err
	}

	// check if memory exists
	count, err := s.dao.Memory.Where(s.dao.Memory.ContentMd5.Eq(dataMd5)).Where(s.dao.Memory.UserId.Eq(int64(userId))).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		mem, err := s.dao.Memory.Where(s.dao.Memory.ContentMd5.Eq(dataMd5)).Where(s.dao.Memory.UserId.Eq(int64(userId))).First()
		if err != nil {
			return 0, err
		}

		return mem.Id, err
	}

	emb, err := s.Embedding.TextEmbedding(ctx, []string{data})
	if err != nil {
		return 0, err
	}

	var mem = &entity.Memory{
		Content:    data,
		ContentMd5: dataMd5,
		Vector:     emb[0],
		UserId:     userId,
	}

	err = s.dao.Memory.Create(mem)
	if err != nil {
		return 0, err
	}

	var entityCols = []entity2.Column{
		entity2.NewColumnFloatVector("vector", s.config.OpenAI.EmbeddingDim, emb),
		entity2.NewColumnVarChar("hash", []string{dataMd5}),
		entity2.NewColumnVarChar("model", []string{s.config.OpenAI.EmbeddingModel}),
		entity2.NewColumnInt64("memory_id", []int64{int64(mem.Id)}),
		entity2.NewColumnInt64("user_id", []int64{int64(userId)}),
	}

	// insert to milvus
	_, err = s.Milvus.Upsert(ctx, s.config.Milvus.Collection, "", entityCols...)
	if err != nil {
		return 0, err
	}

	return mem.Id, err
}

func (s *Service) updateMemory(ctx context.Context, memoryId schema.EntityId, data string) (*entity.Memory, error) {
	dataMd5, err := md5.Md5(data)
	if err != nil {
		return nil, err
	}

	// 检查 memId 是否存在
	count, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.Id.Eq(uint(memoryId))).Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		// 不存在
		return nil, fmt.Errorf("memory id not exists")
	}

	mem, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.Id.Eq(uint(memoryId))).First()
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
	_, err = s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.Id.Eq(uint(memoryId))).Updates(mem)

	if err != nil {
		return nil, err
	}

	// 更新 milvus
	var entityCols = []entity2.Column{
		entity2.NewColumnFloatVector("vector", s.config.OpenAI.EmbeddingDim, embed),
		entity2.NewColumnVarChar("hash", []string{dataMd5}),
		entity2.NewColumnVarChar("model", []string{s.config.OpenAI.EmbeddingModel}),
		entity2.NewColumnInt64("memory_id", []int64{int64(memoryId)}),
		entity2.NewColumnInt64("user_id", []int64{int64(mem.UserId)}),
	}

	// insert to milvus
	_, err = s.Milvus.Upsert(ctx, s.config.Milvus.Collection, "", entityCols...)
	if err != nil {
		return nil, err
	}

	return mem, nil

}

func (s *Service) getMemory(ctx context.Context, memoryId schema.EntityId, data string) error {

	return nil
}

func (s *Service) deleteMemory(ctx context.Context, memoryId schema.EntityId) error {
	// 检查 memId 是否存在
	count, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.Id.Eq(uint(memoryId))).Count()
	if err != nil {
		return err
	}

	if count == 0 {
		// 不存在
		return fmt.Errorf("delete failed, memory id not exists")
	}

	mem, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.Id.Eq(uint(memoryId))).First()
	if err != nil {
		return err
	}

	// milvus delete
	ids := entity2.NewColumnInt64("memory_id", []int64{int64(mem.Id)})
	errDelete := s.Milvus.DeleteByPks(ctx, s.config.Milvus.Collection, "", ids)
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

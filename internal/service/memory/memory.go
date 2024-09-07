package memory

import (
	"context"
	"fmt"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"

	"github.com/bytedance/sonic"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	entity2 "github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/tmc/langchaingo/llms"
)

// 也许服务层不需要关心用户 ID，只需要在控制器层判断

func (s *Service) Add(ctx context.Context, data string, userId schema.UserId) error {
	vec, err := s.Embedding.TextEmbedding(ctx, []string{data})
	if err != nil {
		return err
	}
	var history []llms.MessageContent
	history = append(history, llms.TextParts(llms.ChatMessageTypeSystem, "You are an expert at deducing facts, preferences and memories from unstructured text."))
	history = append(history, llms.TextParts(llms.ChatMessageTypeHuman, s.memoryDeductionPrompt(data)))

	extractedMemories, err := s.OpenAI.GenerateContent(ctx, history, llms.WithMaxTokens(1024))
	if err != nil {
		s.Logger.Sugar.Error("Unable to generate memory response, err: " + err.Error())
		return err
	}
	extractedMemoriesText := extractedMemories.Choices[0].Content

	var filter = fmt.Sprintf("user_id == %s && model == %s", userId, s.config.OpenAI.EmbeddingModel)
	sp, err := entity2.NewIndexAUTOINDEXSearchParam(1)
	if err != nil {
		return err
	}
	vector := entity2.FloatVector(vec[0])
	existingMemories, err := s.Milvus.Search(ctx, s.config.Milvus.MemoryCollection,
		[]string{},
		filter,
		[]string{"memory_id"},
		[]entity2.Vector{vector},
		"vector",
		entity2.L2,
		10,
		sp, client.WithLimit(5))

	var LLMMemories []*LLMMemory

	// get all data
	for _, res := range existingMemories {
		var idColumn *entity2.ColumnInt64
		for _, field := range res.Fields {
			if field.Name() == "memory_id" {
				c, ok := field.(*entity2.ColumnInt64)
				if ok {
					idColumn = c
				}
			}
		}

		if idColumn == nil {
			return fmt.Errorf("id column not found")
		}

		for i := 0; i < res.ResultCount; i++ {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				return err
			}

			mem, err := s.dao.Memory.Where(s.dao.Memory.Id.Eq(uint(id))).First()
			if err != nil {
				return err
			}

			LLMMemories = append(LLMMemories, &LLMMemory{
				ResultId: i,
				ID:       schema.EntityId(id),
				Score:    res.Scores[i],
				Memory:   mem.Content,
			})
		}
	}

	j, err := sonic.MarshalString(LLMMemories)

	prompt := s.updateMemoryPrompt(j, extractedMemoriesText)

	// add tools to llm
	history = []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}
	response, err := s.OpenAI.GenerateContent(ctx, history, llms.WithMaxTokens(1000), llms.WithTools(tools))
	if err != nil {
		return err
	}
	history, err = s.executeToolCalls(ctx, userId, history, response)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetMemories(ctx context.Context, userId schema.UserId) ([]*entity.Memory, error) {
	m, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Where(s.dao.Memory.UserId.Eq(userId.String())).Find()

	if err != nil {
		return nil, err
	}

	return m, err
}
func (s *Service) Exists(ctx context.Context, memoryId uint, userId schema.UserId) (bool, error) {
	i, err := s.dao.WithContext(ctx).Memory.
		Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Where(s.dao.Memory.UserId.Eq(userId.String())).
		Where(s.dao.Memory.Id.Eq(memoryId)).Count()

	return i > 0, err

}

func (s *Service) GetMemory(ctx context.Context, memoryId uint, userId schema.UserId) (*entity.Memory, error) {
	i, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).Where(s.dao.Memory.Id.Eq(memoryId)).Count()
	if err != nil {
		return nil, err
	}

	if i == 0 {
		return nil, consts.ErrMemoryNotFound
	}

	m, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).Where(s.dao.Memory.Id.Eq(memoryId)).Where(s.dao.Memory.UserId.Eq(userId.String())).First()

	if err != nil {
		s.Logger.Sugar.Error("Unable to get memories, err: " + err.Error())
	}

	return m, err
}

func (s *Service) Delete(ctx context.Context, memory *entity.Memory) error {
	// 检查用户是否具有
	m, err := s.dao.WithContext(ctx).Memory.
		Where(s.dao.Memory.Id.Eq(uint(memory.Id))).
		Where(s.dao.Memory.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).First()
	if err != nil {
		return err
	}

	err = s.deleteMemory(ctx, m.Id)
	if err != nil {
		return err
	}

	return nil
}

// Purge 清除用户的所有记忆
func (s *Service) Purge(ctx context.Context, userId schema.UserId) error {
	_, err := s.dao.WithContext(ctx).Memory.Where(s.dao.Memory.UserId.Eq(userId.String())).Delete()

	if err != nil {
		return err
	}

	// milvus delete
	var filter = fmt.Sprintf("user_id == %s", userId)
	errDelete := s.Milvus.Delete(ctx, s.config.Milvus.MemoryCollection, "", filter)
	if errDelete != nil {
		return errDelete
	}

	return nil
}

type LLMMemory struct {
	ID       schema.EntityId `json:"id"`
	Memory   string          `json:"memory"`
	Score    float32         `json:"score"`
	ResultId int             `json:"-"`
}

package memory

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
)

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
	}
	extractedMemoriesText := extractedMemories.Choices[0].Content

	var filter = fmt.Sprintf("user_id == %d", userId)
	sp, err := entity.NewIndexAUTOINDEXSearchParam(1)
	if err != nil {
		return err
	}
	vector := entity.FloatVector(vec[0])
	existingMemories, err := s.Milvus.Search(ctx, s.config.Milvus.Collection,
		[]string{},
		filter,
		[]string{"memory_id"},
		[]entity.Vector{vector},
		"vector",
		entity.L2,
		10,
		sp, client.WithLimit(5))

	var LLMMemories []*LLMMemory

	// get all data
	for _, res := range existingMemories {
		var idColumn *entity.ColumnInt64
		for _, field := range res.Fields {
			if field.Name() == "memory_id" {
				c, ok := field.(*entity.ColumnInt64)
				if ok {
					idColumn = c
				}
			}
		}

		for i := 0; i < res.ResultCount; i++ {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				return err
			}

			mem, err := s.dao.Memory.Where(s.dao.Memory.Id.Eq(uint(id))).First()

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

func (s *Service) Delete(memoryId uint) error {

	return nil
}

type LLMMemory struct {
	ID       schema.EntityId `json:"id"`
	Memory   string          `json:"memory"`
	Score    float32         `json:"score"`
	ResultId int             `json:"-"`
}

package memory

import (
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"rag-new/internal/base/conf"
	logger2 "rag-new/internal/base/logger"
	milvus2 "rag-new/internal/base/milvus"
	"rag-new/internal/base/orm"
	dao2 "rag-new/internal/dao"
	"rag-new/internal/schema"
	"rag-new/internal/service/embedding"
	"testing"
)

func NewService() *Service {
	logger := logger2.NewZapLogger()
	config := conf.ProviderConfig(logger)
	database := orm.NewGORM(config, logger)
	dao := dao2.NewQuery(database)
	emb := embedding.NewEmbedding(config, logger, dao)
	milvus := milvus2.NewMilvus(config)

	return NewMemory(config, logger, emb, milvus, dao)
}
func TestMemory(t *testing.T) {

}

func TestAdd(t *testing.T) {
	memory := NewService()
	ctx := context.Background()

	likes := []string{
		"喜欢喝可乐",
		"喜欢看书",
		"喜欢阅读",
		"不喜欢喝雪碧",
	}

	const userId = 1

	for _, v := range likes {
		err := memory.Add(ctx, v, userId)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
func TestGeneratePrompt(t *testing.T) {
	memoryService := NewService()
	var ctx = context.Background()

	var filter = "user_id == 1"
	sp, _ := entity.NewIndexAUTOINDEXSearchParam(1)

	vec, err := memoryService.Embedding.TextEmbedding(ctx, []string{"可乐"})
	vector := entity.FloatVector(vec[0])

	existingMemories, err := memoryService.Milvus.Search(ctx, memoryService.config.Milvus.Collection,
		[]string{},
		filter,
		[]string{"memory_id"},
		[]entity.Vector{vector},
		"vector",
		entity.L2,
		10,
		sp, client.WithLimit(5))

	if err != nil {
		panic(err)
	}

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

		if idColumn == nil {
			panic("result field not math")
		}
		for i := 0; i < res.ResultCount; i++ {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				panic(err)
			}

			mem, err := memoryService.dao.Memory.Where(memoryService.dao.Memory.Id.Eq(uint(id))).First()

			LLMMemories = append(LLMMemories, &LLMMemory{
				ResultId: i,
				ID:       schema.EntityId(id),
				Score:    res.Scores[i],
				Memory:   mem.Content,
			})
		}
	}

}

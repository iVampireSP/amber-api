package embedding

import (
	"context"
	"rag-new/internal/base/conf"
	logger2 "rag-new/internal/base/logger"
	"rag-new/internal/base/orm"
	dao2 "rag-new/internal/dao"
	"testing"
)

func NewService() *Service {
	logger := logger2.NewZapLogger()
	config := conf.ProviderConfig(logger)
	database := orm.NewGORM(config, logger)
	dao := dao2.NewQuery(database)

	return NewEmbedding(config, logger, dao)
}

func TestEmbedding(t *testing.T) {
	emb := NewService()
	var ctx = context.Background()

	embeddings, err := emb.TextEmbedding(ctx, []string{"Hello world!"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(embeddings)
}

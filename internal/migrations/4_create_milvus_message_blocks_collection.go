package migrations

import (
	"context"
	"database/sql"
	"errors"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/pressly/goose/v3"
	"strconv"
)

func init() {
	goose.AddMigrationContext(Up4createMilvusCollection, Down4createMilvusCollection)
}

func Up4createMilvusCollection(ctx context.Context, _ *sql.Tx) error {
	var field = []*entity.Field{
		{
			Name:       "block_id",
			PrimaryKey: true,
			AutoID:     false,
			DataType:   entity.FieldTypeInt64,
		},
		{
			Name:       "chat_id",
			PrimaryKey: false,
			AutoID:     false,
			DataType:   entity.FieldTypeInt64,
		},
		{
			Name:       "model",
			PrimaryKey: false,
			AutoID:     false,
			DataType:   entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "255",
			},
		},
		{
			Name:       "vector",
			PrimaryKey: false,
			DataType:   entity.FieldTypeFloatVector,
			TypeParams: map[string]string{
				"dim": strconv.Itoa(Config.OpenAI.EmbeddingDim),
			},
		},
	}

	var schema = &entity.Schema{
		CollectionName:     Config.Milvus.MessageBlockCollection,
		Description:        "",
		AutoID:             true,
		Fields:             field,
		EnableDynamicField: true,
	}

	err := Milvus.CreateCollection(ctx, schema, 1)
	if err != nil {
		return errors.Join(errors.New("failed to create collection"), err)
	}

	marisaTrie := entity.NewGenericIndex("idx_model", entity.Trie, map[string]string{})
	err = Milvus.CreateIndex(ctx, Config.Milvus.MessageBlockCollection, "model", marisaTrie, false)
	if err != nil {
		return errors.Join(errors.New("failed to create model index"), err)

	}

	inverted := entity.NewGenericIndex("idx_chat_id", entity.Inverted, map[string]string{})
	err = Milvus.CreateIndex(ctx, Config.Milvus.MessageBlockCollection, "chat_id", inverted, false)
	if err != nil {
		return errors.Join(errors.New("failed to create chat_id index"), err)

	}

	index, err := entity.NewIndexAUTOINDEX(entity.L2)
	if err != nil {
		return errors.Join(errors.New("auto index error"), err)
	}
	err = Milvus.CreateIndex(ctx, Config.Milvus.MessageBlockCollection, "vector", index, false)
	if err != nil {
		return errors.Join(errors.New("failed to create vector index"), err)
	}

	return err
}

func Down4createMilvusCollection(ctx context.Context, _ *sql.Tx) error {
	err := Milvus.DropCollection(ctx, Config.Milvus.MessageBlockCollection)
	return err
}

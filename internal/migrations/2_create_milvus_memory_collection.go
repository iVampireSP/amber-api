package migrations

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up2createMilvusCollection, Down2createMilvusCollection)
}

func Up2createMilvusCollection(ctx context.Context, _ *sql.Tx) error {
	var field = []*entity.Field{
		{
			Name:       "memory_id",
			PrimaryKey: true,
			AutoID:     false,
			DataType:   entity.FieldTypeInt64,
		},
		{
			Name:       "user_id",
			PrimaryKey: false,
			AutoID:     false,
			DataType:   entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "255",
			},
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
			Name:       "hash",
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
		CollectionName:     Config.Milvus.MemoryCollection,
		Description:        "",
		AutoID:             true,
		Fields:             field,
		EnableDynamicField: true,
	}

	err := Milvus.CreateCollection(ctx, schema, 2)
	if err != nil {
		return errors.Join(errors.New("failed to create collection"), err)
	}

	marisaTrie := entity.NewGenericIndex("idx_model", entity.Trie, map[string]string{})
	err = Milvus.CreateIndex(ctx, Config.Milvus.MemoryCollection, "model", marisaTrie, false)
	if err != nil {
		return errors.Join(errors.New("failed to create model index"), err)

	}

	//inverted := entity.NewGenericIndex("idx_user_id", entity.Trie, map[string]string{})
	marisaTrieHash := entity.NewGenericIndex("idx_user_id", entity.Trie, map[string]string{})
	err = Milvus.CreateIndex(ctx, Config.Milvus.MemoryCollection, "user_id", marisaTrieHash, false)
	if err != nil {
		return errors.Join(errors.New("failed to create user_id index"), err)

	}

	marisaTrieHash = entity.NewGenericIndex("idx_hash", entity.Trie, map[string]string{})

	err = Milvus.CreateIndex(ctx, Config.Milvus.MemoryCollection, "hash", marisaTrieHash, false)
	if err != nil {
		return errors.Join(errors.New("failed to create hash index"), err)
	}

	index, err := entity.NewIndexAUTOINDEX(entity.L2)
	if err != nil {
		return errors.Join(errors.New("auto index error"), err)
	}
	err = Milvus.CreateIndex(ctx, Config.Milvus.MemoryCollection, "vector", index, false)
	if err != nil {
		return errors.Join(errors.New("failed to create vector index"), err)
	}

	return err
}

func Down2createMilvusCollection(ctx context.Context, _ *sql.Tx) error {
	err := Milvus.DropCollection(ctx, Config.Milvus.MemoryCollection)
	return err
}

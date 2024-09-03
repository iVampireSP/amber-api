package milvus

import (
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"rag-new/internal/base/conf"
	"strconv"
)

func NewMilvus(config *conf.Config) client.Client {
	c, err := client.NewClient(context.Background(), client.Config{
		Address: config.Milvus.Host + ":" + strconv.Itoa(config.Milvus.Port),
		DBName:  config.Milvus.DBName,
	})

	if err != nil {
		panic(err)
	}

	return c
}

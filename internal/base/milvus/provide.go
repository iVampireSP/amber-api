package milvus

import (
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"strconv"
)

func NewMilvus(config *conf.Config, logger *logger.Logger) client.Client {
	var address = config.Milvus.Host + ":" + strconv.Itoa(config.Milvus.Port)

	logger.Sugar.Infof("Connection to milvus, address=%s, dbname=%s", address, config.Milvus.DBName)

	c, err := client.NewClient(context.Background(), client.Config{
		Address: address,
		DBName:  config.Milvus.DBName,
	})

	if err != nil {
		panic(err)
	}

	return c
}

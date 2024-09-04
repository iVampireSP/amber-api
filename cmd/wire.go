//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	v1 "rag-new/internal/api/v1"
	"rag-new/internal/base"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/base/milvus"
	"rag-new/internal/base/openai"
	"rag-new/internal/base/orm"
	"rag-new/internal/base/redis"
	"rag-new/internal/base/s3"
	"rag-new/internal/base/server"
	"rag-new/internal/batch"
	"rag-new/internal/dao"
	"rag-new/internal/middleware"
	"rag-new/internal/router"
	"rag-new/internal/service"
)

var ProviderSet = wire.NewSet(
	conf.ProviderConfig,
	logger.NewZapLogger,
	milvus.NewMilvus,
	orm.NewGORM,
	dao.NewQuery,
	redis.NewRedis,
	s3.NewS3,
	openai.NewOpenAI,
	middleware.Provider,
	batch.NewBatch,
	service.Provider,
	v1.ProviderApiControllerSet,
	router.ProviderSetRouter,
	server.NewHTTPServer,
	base.NewApplication,
)

func CreateApp() (*base.Application, error) {
	wire.Build(ProviderSet)

	return nil, nil
}

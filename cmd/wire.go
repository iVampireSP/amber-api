//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	v1 "rag-new/internal/api/v1"
	"rag-new/internal/base"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/orm"
	"rag-new/internal/base/server"
	"rag-new/internal/logger"
	"rag-new/internal/middleware"
	"rag-new/internal/router"
	"rag-new/internal/service"
)

var ProviderSet = wire.NewSet(
	conf.ProviderConfig,
	logger.NewZapLogger,
	orm.NewXORM,
	middleware.Provider,
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

//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	v1 "rag-new/internal/api/v1"
	"rag-new/internal/base"
	"rag-new/internal/base/conf"
	"rag-new/internal/logger"
	"rag-new/internal/orm"
	"rag-new/internal/router"
	"rag-new/internal/server"
	"rag-new/internal/services"
)

var ProviderSet = wire.NewSet(
	conf.ProviderConfig,
	logger.NewZapLogger,
	orm.NewXORM,
	services.Provider,
	v1.ProviderApiControllerSet,
	router.ProviderSetRouter,
	server.NewHTTPServer,
	base.NewApplication,
)

func CreateApp() (*base.Application, error) {
	wire.Build(ProviderSet)

	return nil, nil
}

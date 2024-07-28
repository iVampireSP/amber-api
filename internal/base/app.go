package base

import (
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/base/server"
	"rag-new/internal/middleware"
	"rag-new/internal/service"
	"xorm.io/xorm"
)

type Application struct {
	Config     *conf.Config
	HttpServer *server.HttpServer
	Logger     *logger.Logger
	X          *xorm.Engine
	Service    *service.Service
	Middleware *middleware.Middleware
}

func NewApplication(
	config *conf.Config,
	httpServer *server.HttpServer,
	logger *logger.Logger,
	x *xorm.Engine,
	services *service.Service,
	middleware *middleware.Middleware,
) *Application {
	return &Application{
		Config:     config,
		HttpServer: httpServer,
		Logger:     logger,
		X:          x,
		Service:    services,
		Middleware: middleware,
	}
}

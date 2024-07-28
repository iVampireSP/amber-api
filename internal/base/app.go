package base

import (
	"github.com/gin-gonic/gin"
	"rag-new/internal/base/conf"
	"rag-new/internal/logger"
	"rag-new/internal/middleware"
	"rag-new/internal/services"
	"xorm.io/xorm"
)

type Application struct {
	Config     *conf.Config
	Gin        *gin.Engine
	Logger     *logger.Logger
	X          *xorm.Engine
	Service    *services.Service
	Middleware *middleware.Middleware
}

func NewApplication(
	config *conf.Config,
	gin *gin.Engine,
	logger *logger.Logger,
	x *xorm.Engine,
	services *services.Service,
	middleware *middleware.Middleware,
) *Application {
	return &Application{
		Config:     config,
		Gin:        gin,
		Logger:     logger,
		X:          x,
		Service:    services,
		Middleware: middleware,
	}
}

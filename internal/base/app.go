package base

import (
	"github.com/gin-gonic/gin"
	"rag-new/internal/base/conf"
	"rag-new/internal/logger"
	"xorm.io/xorm"
)

type Application struct {
	Config *conf.Config
	Gin    *gin.Engine
	Logger *logger.Logger
	X      *xorm.Engine
}

func NewApplication(
	config *conf.Config,
	gin *gin.Engine,
	logger *logger.Logger,
	x *xorm.Engine,
) *Application {
	return &Application{
		Config: config,
		Gin:    gin,
		Logger: logger,
		X:      x,
	}
}

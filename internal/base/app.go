package base

import (
	"github.com/gin-gonic/gin"
	"rag-new/internal/base/conf"
	"rag-new/internal/logger"
)

type Application struct {
	Config *conf.Config
	Gin    *gin.Engine
	Logger *logger.Logger
}

func NewApplication(
	config *conf.Config,
	gin *gin.Engine,
	logger *logger.Logger,
) *Application {
	return &Application{
		Config: config,
		Gin:    gin,
		Logger: logger,
	}
}

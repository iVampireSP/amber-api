package base

import (
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/base/s3"
	"rag-new/internal/base/server"
	"rag-new/internal/batch"
	"rag-new/internal/dao"
	"rag-new/internal/middleware"
	"rag-new/internal/service"
)

type Application struct {
	Config     *conf.Config
	HttpServer *server.HttpServer
	Logger     *logger.Logger
	GORM       *gorm.DB
	DAO        *dao.Query
	Service    *service.Service
	Middleware *middleware.Middleware
	Redis      *redis.Client
	Batch      *batch.Batch
	S3         *s3.S3
	OpenAI     *openai.Client
}

func NewApplication(
	config *conf.Config,
	httpServer *server.HttpServer,
	logger *logger.Logger,
	services *service.Service,
	middleware *middleware.Middleware,
	redis *redis.Client,
	batch *batch.Batch,
	S3 *s3.S3,
	GORM *gorm.DB,
	DAO *dao.Query,
	OpenAI *openai.Client,
) *Application {
	return &Application{
		Config:     config,
		HttpServer: httpServer,
		Logger:     logger,
		Service:    services,
		Middleware: middleware,
		Redis:      redis,
		Batch:      batch,
		S3:         S3,
		GORM:       GORM,
		DAO:        DAO,
		OpenAI:     OpenAI,
	}
}

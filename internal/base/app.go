package base

import (
	milvusClient "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/base/redis"
	"rag-new/internal/base/s3"
	"rag-new/internal/base/server"
	"rag-new/internal/batch"
	"rag-new/internal/dao"
	"rag-new/internal/middleware"
	"rag-new/internal/service"
	"rag-new/internal/service/embedding"
)

type Application struct {
	Config     *conf.Config
	HttpServer *server.HttpServer
	Logger     *logger.Logger
	GORM       *gorm.DB
	DAO        *dao.Query
	Service    *service.Service
	Middleware *middleware.Middleware
	Redis      *redis.Redis
	Batch      *batch.Batch
	S3         *s3.S3
	OpenAI     *openai.Client
	Milvus     milvusClient.Client
	Embedding  *embedding.Service
}

func NewApplication(
	config *conf.Config,
	httpServer *server.HttpServer,
	logger *logger.Logger,
	services *service.Service,
	middleware *middleware.Middleware,
	redis *redis.Redis,
	batch *batch.Batch,
	S3 *s3.S3,
	GORM *gorm.DB,
	DAO *dao.Query,
	OpenAI *openai.Client,
	Milvus milvusClient.Client,
	embedding *embedding.Service,
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
		Milvus:     Milvus,
		Embedding:  embedding,
	}
}

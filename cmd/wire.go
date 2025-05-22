//go:build wireinject
// +build wireinject

package cmd

import (
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
	"rag-new/internal/message"
	"rag-new/internal/middleware"
	"rag-new/internal/router"
	"rag-new/internal/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// Base
	conf.ProviderConfig,
	logger.NewZapLogger,
	message.NewMessage,
	milvus.NewMilvus,
	orm.NewGORM,
	dao.NewQuery,
	redis.NewRedis,
	s3.NewS3,
	openai.NewOpenAI,
	batch.NewBatch,

	// Services
	// 全部到 internal/service/provider.go 注册
	// jwks.NewJWKS,
	// auth.NewAuthService,
	// embedding.NewService,
	// memory.NewMemory,
	// chat_message.NewService,
	// chat.NewService,
	// tool.NewService,
	// assistant.NewService,
	// builtin_tool.NewService,
	// llm.NewLLM,
	// message_block.NewService,
	// file.NewService,
	// stream.NewService,
	// library.NewService,
	// token_usage.NewService,
	// account.NewService,
	// unsettled_token.NewService,
	// text_classification.NewService,
	service.Provider,

	// Middleware
	middleware.NewMiddleware,

	// Http Controller
	v1.NewUserController,
	v1.NewAuthController,
	v1.NewToolController,
	v1.NewAssistantController,
	v1.NewChatController,
	v1.NewFileController,
	v1.NewMemoryController,
	v1.NewLibraryController,
	v1.NewUsageController,

	// Router
	router.NewApiRoute,
	router.NewSwaggerRoute,

	// Application
	server.NewHTTPServer,
	base.NewApplication,
)

func CreateApp() (*base.Application, error) {
	wire.Build(ProviderSet)

	return nil, nil
}

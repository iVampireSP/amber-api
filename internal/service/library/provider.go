package library

import (
	milvusClient "github.com/milvus-io/milvus-sdk-go/v2/client"
	"rag-new/internal/base/conf"
	"rag-new/internal/dao"
	"rag-new/internal/service/embedding"
	"rag-new/internal/service/file"
)

type Service struct {
	config      *conf.Config
	dao         *dao.Query
	milvus      milvusClient.Client
	fileService *file.Service
	embedding   *embedding.Service
}

func NewService(config *conf.Config, dao *dao.Query, milvus milvusClient.Client, fileService *file.Service, embedding *embedding.Service) *Service {
	return &Service{
		config:      config,
		dao:         dao,
		milvus:      milvus,
		fileService: fileService,
		embedding:   embedding,
	}
}

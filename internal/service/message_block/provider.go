package message_block

import (
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"rag-new/internal/base/conf"
	"rag-new/internal/dao"
	"rag-new/internal/service/embedding"
)

type Service struct {
	dao       *dao.Query
	config    *conf.Config
	embedding *embedding.Service
	milvus    client.Client
}

func NewService(dao *dao.Query,
	config *conf.Config,
	embedding *embedding.Service,
	milvus client.Client,
) *Service {
	return &Service{dao, config, embedding, milvus}
}

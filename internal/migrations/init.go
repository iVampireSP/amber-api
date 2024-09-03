package migrations

import (
	"embed"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"rag-new/internal/base/conf"
)

//go:embed *.sql
var MigrationFS embed.FS

var Config *conf.Config
var Milvus client.Client

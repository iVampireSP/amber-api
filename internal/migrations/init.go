package migrations

import (
	"embed"
	"rag-new/internal/base/conf"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

//go:embed *.sql
var MigrationFS embed.FS

var Config *conf.Config
var Milvus client.Client

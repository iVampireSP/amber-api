package builtin_tool

import (
	"context"
	"rag-new/internal/base/conf"
	logger2 "rag-new/internal/base/logger"
	"rag-new/internal/base/s3"
	"rag-new/internal/schema"
	"rag-new/internal/service/file"
	"testing"
)

func TestReader(t *testing.T) {
	var logger = logger2.NewZapLogger()
	var config = conf.ProviderConfig(logger)
	var s3Service = s3.NewS3(config)
	var fileService = file.NewService(s3Service, config, nil, nil, logger)
	var s = NewService(config, logger, fileService, nil)

	var ctx = context.Background()

	response, err := s.ReadUrl(ctx, schema.FunctionCallArguments{
		"url": "https://www.baidu.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(response)

	response, err = s.ReadUrl(ctx, schema.FunctionCallArguments{
		"url": "https://exmaple_not_exists.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(response)
}

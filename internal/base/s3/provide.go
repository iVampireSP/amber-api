package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"rag-new/internal/base/conf"
)

func NewS3(config *conf.Config) *minio.Client {
	minioClient, err := minio.New(config.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3.AccessKey, config.S3.SecretKey, ""),
		Secure: config.S3.UseSSL,
	})

	if err != nil {
		panic(err)
	}

	return minioClient
}

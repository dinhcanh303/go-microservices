package minio

import (
	"context"
	"io"
	"mime/multipart"

	minioV7 "github.com/minio/minio-go/v7"
)

type MinioService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, buffer io.Reader) (*minioV7.UploadInfo, error)
	DeleteFile(ctx context.Context, fileNames []string) (bool, error)
}

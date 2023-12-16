package minio

import (
	"context"
	"io"
	"mime/multipart"

	minioV7 "github.com/minio/minio-go/v7"
)

type MinioService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, buffer io.Reader, location string) (*minioV7.UploadInfo, *FileInfo, error)
	DeleteFile(ctx context.Context, objectName string) (bool, error)
}

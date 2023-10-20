package minio

import (
	"io"
	"mime/multipart"

	minioV7 "github.com/minio/minio-go/v7"
)

type MinioUpload interface {
	UploadFile(file *multipart.FileHeader, buffer io.Reader) (*minioV7.UploadInfo, error)
}

package minio

import (
	"context"
	"io"
	"log"
	"log/slog"
	"mime/multipart"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	minioV7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

type minio struct {
	cf *configs.Minio
}

func NewMinio(cf *configs.Minio) MinioService {
	return &minio{
		cf: cf,
	}
}

// DeleteFile implements MinioService.
func (m *minio) DeleteFile(ctx context.Context, fileNames []string) (bool, error) {
	endpoint := m.cf.EndPoint
	accessKeyID := m.cf.AccessKeyID
	secretAccessKey := m.cf.SecretAccessKey
	useSSL := m.cf.UseSSL
	minioClient, err := minioV7.New(endpoint, &minioV7.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return false, errors.Wrap(err, "minio.DeleteFile failed")
	}
	opts := minioV7.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	for _, fileName := range fileNames {
		err := minioClient.RemoveObject(ctx, m.cf.BucketName, fileName, opts)
		if err != nil {
			return false, errors.Wrap(err, "minio.DeleteFile failed")
		}
	}
	return true, nil
}

// UploadFile implements MinioUpload.
func (m *minio) UploadFile(ctx context.Context, file *multipart.FileHeader, buffer io.Reader) (*minioV7.UploadInfo, error) {
	slog.Info("MINIO::Upload File")
	mdf.Pla
	endpoint := m.cf.EndPoint
	accessKeyID := m.cf.AccessKeyID
	secretAccessKey := m.cf.SecretAccessKey
	// useSSL := m.cf.UseSSL
	region := m.cf.Region
	bucketName := m.cf.BucketName
	// rootForlder := m.cf.RootFolder
	// Initialize minio client object.
	minioClient, err := minioV7.New(endpoint, &minioV7.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
		Region: region,
	})
	if err != nil {
		log.Fatalln(err)
	}
	objectName := file.Filename
	contentType := file.Header.Get("Content-Type")
	fileSize := file.Size
	slog.Info("FILE::", objectName, contentType, fileSize)
	info, err := minioClient.PutObject(ctx, bucketName, objectName, buffer, fileSize, minioV7.PutObjectOptions{
		ContentType: contentType,
	})
	slog.Info("MINIO FILE ERROR::", err)
	slog.Info("MINIO FILE::", info)
	if err != nil {
		log.Fatalln(err)
	}
	slog.Info("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return &info, nil
}

var _ MinioService = (*minio)(nil)

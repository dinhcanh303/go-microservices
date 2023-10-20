package minio

import (
	"context"
	"io"
	"log"
	"mime/multipart"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	minioV7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

type minio struct {
	cf configs.Minio
}

// UploadFile implements MinioUpload.
func (m *minio) UploadFile(file *multipart.FileHeader, buffer io.Reader) (*minioV7.UploadInfo, error) {
	ctx := context.Background()
	endpoint := m.cf.EndPoint
	accessKeyID := m.cf.AccessKeyID
	secretAccessKey := m.cf.SecretAccessKey
	useSSL := m.cf.UseSSL

	// Initialize minio client object.
	minioClient, err := minioV7.New(endpoint, &minioV7.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	// Make a new bucket called.
	bucketName := m.cf.BucketName
	location := m.cf.Location
	err = minioClient.MakeBucket(ctx, bucketName, minioV7.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return nil, errors.Wrap(err, "We already own ::")
		} else {
			return nil, errors.Wrap(err, "Error:")
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the zip file with FPutObject
	// info, err := minioClient.FPutObject(ctx, bucketName, objectName,
	// 	filePath, minioV7.PutObjectOptions{ContentType: contentType})
	objectName := file.Filename
	contentType := file.Header.Get("Content-Type")
	fileSize := file.Size
	info, err := minioClient.PutObject(ctx, bucketName, objectName, buffer, fileSize, minioV7.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return &info, nil
}

var _ MinioUpload = (*minio)(nil)

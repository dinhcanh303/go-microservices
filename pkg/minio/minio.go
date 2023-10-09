package minio

import (
	"context"
	"log"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	minioV7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minio struct {
	cf configs.Minio
}

// UploadFile implements MinioUpload.
func (m *minio) UploadFile() {
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
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the zip file
	objectName := "golden-oldies.zip"
	filePath := "/tmp/golden-oldies.zip"
	contentType := "application/zip"

	// Upload the zip file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minioV7.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

var _ MinioUpload = (*minio)(nil)

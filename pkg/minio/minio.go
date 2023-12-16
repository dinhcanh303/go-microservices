package minio

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	minioV7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type minio struct {
	cf *configs.Minio
}
type FileInfo struct {
	FileName     string `json:"filename"`
	Extension    string `json:"extension"`
	MimeType     string `json:"mime_type"`
	Folder       string `json:"folder"`
	Url          string `json:"url"`
	UrlThumbnail string `json:"url_thumbnail"`
}

func NewMinio(cf *configs.Minio) MinioService {
	return &minio{
		cf: cf,
	}
}

// DeleteFile implements MinioService.
func (m *minio) DeleteFile(ctx context.Context, fileName string) (bool, error) {
	config := m.cf
	client, _ := minioClient(config)
	opts := minioV7.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err := client.RemoveObject(ctx, m.cf.BucketName, fileName, opts)
	if err != nil {
		return false, status.Error(codes.AlreadyExists, "RemoveObject failed")
	}
	return true, nil
}

// UploadFile implements MinioUpload.
func (m *minio) UploadFile(ctx context.Context, file *multipart.FileHeader, buffer io.Reader, location string) (*minioV7.UploadInfo, *FileInfo, error) {
	slog.Info("MINIO::Upload File")
	config := m.cf
	client, _ := minioClient(config)
	folder := locationFolderSaveFile(config.RootFolder, location)
	objectName := folder + file.Filename
	contentType := file.Header.Get("Content-Type")
	fileSize := file.Size
	info, err := client.PutObject(ctx, config.BucketName, objectName, buffer, fileSize, minioV7.PutObjectOptions{
		ContentType: contentType,
	})
	urlFile := getUrlFile(config.EndPoint, config.BucketName, objectName)
	fileInfo := &FileInfo{
		FileName:     file.Filename,
		Folder:       folder,
		Extension:    getExtensionFile(file.Filename),
		MimeType:     contentType,
		Url:          urlFile,
		UrlThumbnail: urlFile,
	}
	if err != nil {
		log.Fatalln(err)
	}
	slog.Info("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return &info, fileInfo, nil
}
func getExtensionFile(fileName string) string {
	extension := filepath.Ext(fileName)
	return extension[1:]
}
func minioClient(config *configs.Minio) (*minioV7.Client, error) {
	endpoint := config.EndPoint
	accessKeyID := config.AccessKeyID
	secretAccessKey := config.SecretAccessKey
	useSSL := config.UseSSL
	minioClient, err := minioV7.New(endpoint, &minioV7.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, "Client Minio failed")
	}
	return minioClient, nil
}
func getUrlFile(endpoint string, bucketName string, objectName string) string {
	return fmt.Sprintf("https://%s/%s/%s", endpoint, bucketName, objectName)
}
func locationFolderSaveFile(rootFolder string, location string) string {
	year, month, _ := time.Now().Date()
	if location == "" {
		return fmt.Sprintf("%s/%d/%02d/", rootFolder, year, int(month))
	}
	return fmt.Sprintf("%s/%s/%d/%02d/", rootFolder, location, year, int(month))
}

var _ MinioService = (*minio)(nil)

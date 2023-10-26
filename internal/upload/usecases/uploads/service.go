package uploads

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/pkg/minio"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type uploadService struct {
	repo  AttachmentRepo
	minio minio.MinioService
}

var _ UseCase = (*uploadService)(nil)

var UseCaseSet = wire.NewSet(NewUploadService)

func NewUploadService(repo AttachmentRepo, minio minio.MinioService) UseCase {
	return &uploadService{
		repo:  repo,
		minio: minio,
	}
}

// DeleteAttachment implements UseCase.
func (s *uploadService) DeleteAttachment(ctx context.Context, attachmentId uuid.UUID) (bool, error) {
	attachment, err := s.repo.Get(ctx, attachmentId)
	if err != nil {
		return false, errors.Wrap(err, "uploadService.DeleteAttachment failed")
	}
	fileNames := make([]string, 1)
	fileNames = append(fileNames, attachment.FileName)
	deletedFile, err := s.minio.DeleteFile(ctx, fileNames)
	if err != nil {
		return false, errors.Wrap(err, "minioService.DeleteFile failed")
	}
	if !deletedFile {
		return false, errors.Wrap(err, "minioService.DeleteFile failed")
	}
	deleted, err := s.repo.Delete(ctx, attachmentId)
	if err != nil {
		return false, errors.Wrap(err, "uploadService.DeleteAttachment failed")
	}
	return deleted, nil
}

// GetAttachment implements UseCase.
func (s *uploadService) GetAttachment(ctx context.Context, attachmentId uuid.UUID) (*domain.Attachment, error) {
	attachment, err := s.repo.Get(ctx, attachmentId)
	if err != nil {
		return nil, errors.Wrap(err, "uploadService.GetAttachment failed")
	}
	return attachment, nil
}

// UpdateAttachment implements UseCase.
func (s *uploadService) UpdateAttachment(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error) {
	result, err := s.repo.Update(ctx, &domain.Attachment{
		AttachableType: attachment.AttachableType,
		AttachableID:   attachment.AttachableID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "uploadService.UpdateAttachment failed")
	}
	return result, nil
}

// UploadFile implements UseCase.
func (s *uploadService) UploadFile(echoCtx echo.Context) ([]*domain.Attachment, error) {
	slog.Info("Service: UploadFile")
	ctx := context.Background()
	form, err := echoCtx.MultipartForm()
	if err != nil {
		return nil, errors.Wrap(err, "Get Upload Form Error")
	}
	slog.Info("FORMDATA::", form)
	files := form.File["files"]
	results := make([]*domain.Attachment, 1)
	for _, file := range files {
		buffer, err := file.Open()
		if err != nil {
			return nil, errors.Wrap(err, "Open Upload file buffer failed")
		}
		defer buffer.Close()
		slog.Info("FILE::", &file)
		info, err := s.minio.UploadFile(ctx, file, buffer)
		slog.Info("INFO::", info)
		if err != nil {
			return nil, errors.Wrap(err, "Upload file failed:")
		}
		model := domain.NewAttachment(uuid.New(), info.Key, "", "", info.Bucket, info.VersionID, info.Location, info.Location)
		attachment, err := s.repo.Create(ctx, model)
		if err != nil {
			return nil, errors.Wrap(err, "Create Attachment failed")
		}
		results = append(results, attachment)

	}
	slog.Info("Data::", results)
	return results, nil
}

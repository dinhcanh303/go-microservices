package uploads

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/pkg/minio"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type uploadService struct {
	repo  AttachmentRepo
	minio minio.MinioService
}

func NewUploadService(repo AttachmentRepo, minio minio.MinioService) UseCase {
	return &uploadService{
		repo:  repo,
		minio: minio,
	}
}

// DeleteAttachment implements UseCase.
func (s *uploadService) DeleteAttachment(ctx context.Context, attachmentId uuid.UUID) (bool, error) {
	panic("unimplemented")
}

// GetAttachment implements UseCase.
func (s *uploadService) GetAttachment(ctx context.Context, attachmentId uuid.UUID) (*domain.Attachment, error) {
	panic("unimplemented")
}

// UpdateAttachment implements UseCase.
func (s *uploadService) UpdateAttachment(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error) {
	panic("unimplemented")
}

// UploadFile implements UseCase.
func (s *uploadService) UploadFile(ctx echo.Context) ([]*string, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, errors.Wrap(err, "Get Upload Form Error")
	}
	files := form.File["files"]
	for _, file := range files {
		buffer, err := file.Open()
		if err != nil {
			return nil, errors.Wrap(err, "Open Upload file buffer failed")
		}
		defer buffer.Close()
		infoFile, err := s.minio.UploadFile(file, buffer)
		if err != nil {
			slog.Warn("Upload file failed: %v", err)
		}
		s.repo.Create()
	}
	return nil, nil
}

// var _ UseCase = (*uploadService)(nil)

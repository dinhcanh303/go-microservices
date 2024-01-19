package uploads

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/pkg/minio"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
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

// UpdateAttachmentsByIds implements UseCase.
func (s *uploadService) UpdateAttachmentsByIds(ctx context.Context, attachmentIds []uuid.UUID, attachment *domain.Attachment) ([]*domain.Attachment, error) {
	results, err := s.repo.UpdateByIds(ctx, attachmentIds, &domain.Attachment{
		AttachableType: attachment.AttachableType,
		AttachableID:   attachment.AttachableID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "uploadService.UpdateAttachmentsByIds failed")
	}
	return results, nil
}

var _ UseCase = (*uploadService)(nil)

var UseCaseSet = wire.NewSet(NewUploadService)

func NewUploadService(repo AttachmentRepo, minio minio.MinioService) UseCase {
	return &uploadService{
		repo:  repo,
		minio: minio,
	}
}

// DeleteAttachmentsByIds implements UseCase.
func (s *uploadService) DeleteAttachmentsByIds(ctx context.Context, attachmentIds []uuid.UUID) (bool, error) {
	attachments, err := s.repo.GetByIds(ctx, attachmentIds)
	if err != nil {
		return false, errors.Wrap(err, "uploadService.DeleteAttachment failed")
	}
	for _, attachment := range attachments {
		_, err := s.minio.DeleteFile(ctx, attachment.Folder+attachment.FileName)
		if err != nil {
			slog.Warn("Minio.DeleteFile failed", err)
		}
	}
	deleted, err := s.repo.DeleteByIds(ctx, attachmentIds)
	if err != nil {
		return false, errors.Wrap(err, "uploadService.DeleteAttachment failed")
	}
	return deleted, nil
}

// DeleteAttachment implements UseCase.
func (s *uploadService) DeleteAttachment(ctx context.Context, attachmentId uuid.UUID) (bool, error) {
	attachment, err := s.repo.Get(ctx, attachmentId)
	if err != nil {
		return false, errors.Wrap(err, "uploadService.DeleteAttachment failed")
	}
	deletedFile, err := s.minio.DeleteFile(ctx, attachment.Folder+attachment.FileName)
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

// GetAttachmentByIds implements UseCase.
func (s *uploadService) GetAttachmentByIds(ctx context.Context, attachmentIds []uuid.UUID) ([]*domain.Attachment, error) {
	attachments, err := s.repo.GetByIds(ctx, attachmentIds)
	if err != nil {
		return nil, errors.Wrap(err, "uploadService.GetAttachment failed")
	}
	return attachments, nil
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
		ID:             attachment.ID,
		AttachableType: attachment.AttachableType,
		AttachableID:   attachment.AttachableID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "uploadService.UpdateAttachment failed")
	}
	return result, nil
}

// UploadFile implements UseCase.
func (s *uploadService) UploadFile(echoCtx echo.Context, location string) ([]*domain.Attachment, error) {
	slog.Info("Service: UploadFile")
	ctx := context.Background()
	form, err := echoCtx.MultipartForm()
	slog.Info("Form", form)
	if err != nil {
		return nil, errors.Wrap(err, "Get Upload Form Error")
	}
	user, err := utils.ExtractHeaderUser(echoCtx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Header User failed")
	}
	files := form.File["files"]
	results := make([]*domain.Attachment, 0)
	for _, file := range files {
		buffer, err := file.Open()
		if err != nil {
			return nil, errors.Wrap(err, "Open Upload file buffer failed")
		}
		defer buffer.Close()
		slog.Info("FILE::", &file)
		_, fileInfo, err := s.minio.UploadFile(ctx, file, buffer, location)
		if err != nil {
			return nil, errors.Wrap(err, "Upload file failed:")
		}
		model := domain.NewAttachment(user.ID, fileInfo.FileName, fileInfo.Extension,
			fileInfo.MimeType, fileInfo.Folder, fileInfo.Url, fileInfo.UrlThumbnail)
		attachment, err := s.repo.Create(ctx, model)
		if err != nil {
			return nil, errors.Wrap(err, "Create Attachment failed")
		}
		results = append(results, attachment)

	}
	return results, nil
}

// func resizeImage(imgIn io.Reader, width, height uint) (io.Reader, error) {
// 	img, _, err := image.Decode(imgIn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
// 	var outBuffer bytes.Buffer
// 	// if err :=
// }

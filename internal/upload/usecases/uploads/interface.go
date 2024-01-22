package uploads

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	AttachmentRepo interface {
		Get(ctx context.Context, attachmentId uuid.UUID) (*domain.Attachment, error)
		GetByIds(ctx context.Context, attachmentIds []uuid.UUID) ([]*domain.Attachment, error)
		Create(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error)
		Update(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error)
		UpdateByIds(ctx context.Context, attachmentIds []uuid.UUID, attachment *domain.Attachment) ([]*domain.Attachment, error)
		Delete(ctx context.Context, attachmentId uuid.UUID) (bool, error)
		DeleteByIds(ctx context.Context, attachmentIds []uuid.UUID) (bool, error)
		GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domain.Attachment, error)
		GetAttachmentsByOptional(ctx context.Context, attachment *domain.Attachment) ([]*domain.Attachment, error)
		GetLastAttachmentByType(ctx context.Context, attachableType string, attachableId uuid.UUID) (*domain.Attachment, error)
	}
	UseCase interface {
		GetAttachment(ctx context.Context, attachmentId uuid.UUID) (*domain.Attachment, error)
		GetAttachmentByIds(ctx context.Context, attachmentIds []uuid.UUID) ([]*domain.Attachment, error)
		UpdateAttachment(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error)
		UpdateAttachmentsByIds(ctx context.Context, attachmentIds []uuid.UUID, attachment *domain.Attachment) ([]*domain.Attachment, error)
		DeleteAttachment(ctx context.Context, attachmentId uuid.UUID) (bool, error)
		DeleteAttachmentsByIds(ctx context.Context, attachmentIds []uuid.UUID) (bool, error)
		UploadFile(ctx echo.Context, location string) ([]*domain.Attachment, error)
	}
	UseCaseGRPC interface {
		GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domain.Attachment, error)
		GetLastAttachmentByType(ctx context.Context, attachableType string, attachableId uuid.UUID) (*domain.Attachment, error)
		GetAttachmentsByOptional(ctx context.Context, attachment *domain.Attachment) ([]*domain.Attachment, error)
	}
)

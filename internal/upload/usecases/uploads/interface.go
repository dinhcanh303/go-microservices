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
		Create(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error)
		Update(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error)
		Delete(ctx context.Context, attachmentId uuid.UUID) (bool, error)
	}
	UseCase interface {
		GetAttachment(ctx context.Context, attachmentId uuid.UUID) (*domain.Attachment, error)
		UpdateAttachment(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error)
		DeleteAttachment(ctx context.Context, attachmentId uuid.UUID) (bool, error)
		UploadFile(ctx echo.Context) ([]*domain.Attachment, error)
	}
)

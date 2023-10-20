package attachments

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type service struct {
	repo AttachmentRepo
}

// CreateAttachment implements UseCase.
func (*service) CreateAttachment(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error) {
	panic("unimplemented")
}

// DeleteAttachment implements UseCase.
func (*service) DeleteAttachment(ctx context.Context, attachmentId uuid.UUID) (bool, error) {
	panic("unimplemented")
}

// GetAttachment implements UseCase.
func (*service) GetAttachment(ctx context.Context, attachmentId uuid.UUID) (*domain.Attachment, error) {
	panic("unimplemented")
}

// UpdateAttachment implements UseCase.
func (*service) UpdateAttachment(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error) {
	panic("unimplemented")
}

var _ UseCase = (*service)(nil)

func NewService(repo AttachmentRepo) *service {
	return &service{
		repo: repo,
	}
}

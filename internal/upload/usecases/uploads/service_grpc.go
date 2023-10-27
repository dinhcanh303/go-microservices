package uploads

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type uploadGRPCService struct {
	repo AttachmentRepo
}

// GetAttachmentsByType implements UseCaseGRPC.
func (s *uploadGRPCService) GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domain.Attachment, error) {
	attachments, err := s.repo.GetAttachmentsByType(ctx, attachableType, attachableId)
	if err != nil {
		return nil, errors.Wrap(err, "uploadService.GetAttachmentsByType failed")
	}
	return attachments, nil
}

var _ UseCaseGRPC = (*uploadGRPCService)(nil)

var UseCaseGRPCSet = wire.NewSet(NewUploadGRPCService)

func NewUploadGRPCService(repo AttachmentRepo) UseCaseGRPC {
	return &uploadGRPCService{
		repo: repo,
	}
}

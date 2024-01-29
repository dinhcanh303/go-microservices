package uploads

import (
	"context"
	"log/slog"
	"strings"

	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type uploadGRPCService struct {
	repo  AttachmentRepo
	redis redis.RedisEngine
}

var _ UseCaseGRPC = (*uploadGRPCService)(nil)

var UseCaseGRPCSet = wire.NewSet(NewUploadGRPCService)

func NewUploadGRPCService(
	repo AttachmentRepo,
	redis redis.RedisEngine,
) UseCaseGRPC {
	return &uploadGRPCService{
		repo:  repo,
		redis: redis,
	}
}

// GetAttachmentsByOptional implements UseCaseGRPC.
func (s *uploadGRPCService) GetAttachmentsByOptional(ctx context.Context, attachment *domain.Attachment) ([]*domain.Attachment, error) {
	var attachments []*domain.Attachment
	keyCache := constant.CACHE_SV_UPLOAD_ATTACHMENTS +
		attachment.UserID.String() + "_" + strings.ToLower(attachment.AttachableType) +
		"_" + attachment.EntityUploadID + "_" + attachment.MimeType
	err := utils.HandleHitCache(attachments, s.redis, keyCache)
	if err != nil {
		attachments, err = s.repo.GetAttachmentsByOptional(ctx, attachment)
		if err != nil {
			return nil, errors.Wrap(err, "uploadService.GetAttachmentsByOptional failed")
		}
		err = s.redis.Set(keyCache, attachments, 0)
		if err != nil {
			slog.Error("set cache attachments failed")
		}
	}
	return attachments, nil
}

// GetLastAttachmentByType implements UseCaseGRPC.
func (s *uploadGRPCService) GetLastAttachmentByType(ctx context.Context, attachableType string, attachableId uuid.UUID) (*domain.Attachment, error) {
	var attachments []*domain.Attachment
	keyCache := constant.CACHE_SV_UPLOAD_ATTACHMENTS + strings.ToLower(attachableType) +
		"_" + attachableId.String()
	err := utils.HandleHitCache(attachments, s.redis, keyCache)
	if err != nil {
		attachments, err = s.repo.GetAttachmentsByType(ctx, attachableType, attachableId)
		if err != nil {
			return nil, errors.Wrap(err, "uploadService.GetAttachmentsByType failed")
		}
		err = s.redis.Set(keyCache, attachments, 0)
		if err != nil {
			slog.Error("set cache attachments failed")
		}
	}
	if len(attachments) > 0 {
		return attachments[len(attachments)-1], nil
	}
	return nil, nil
}

// GetAttachmentsByType implements UseCaseGRPC.
func (s *uploadGRPCService) GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domain.Attachment, error) {
	var attachments []*domain.Attachment
	keyCache := constant.CACHE_SV_UPLOAD_ATTACHMENTS + strings.ToLower(attachableType) +
		"_" + attachableId.String()
	err := utils.HandleHitCache(attachments, s.redis, keyCache)
	if err != nil {
		attachments, err = s.repo.GetAttachmentsByType(ctx, attachableType, attachableId)
		if err != nil {
			return nil, errors.Wrap(err, "uploadService.GetAttachmentsByType failed")
		}
		err = s.redis.Set(keyCache, attachments, 0)
		if err != nil {
			slog.Error("set cache attachments failed")
		}
	}
	return attachments, nil
}

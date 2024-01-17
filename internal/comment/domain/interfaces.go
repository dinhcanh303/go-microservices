package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
)

type (
	LikeDomainService interface {
		GetLikesInfoByCommentID(ctx context.Context, commentId uuid.UUID, userId uuid.UUID) (*domain.LikesInfo, error)
	}
	UploadDomainService interface {
		GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domainUpload.Attachment, error)
	}
	AuthDomainService interface {
		GetProfile(ctx context.Context, id uuid.UUID) (*gen.GetProfileResponse, error)
	}
)

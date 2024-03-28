package domain

import (
	"context"

	v1u "github.com/dinhcanh303/go-microservices/api/auth/v1"
	v1p "github.com/dinhcanh303/go-microservices/api/post/v1"
	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
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
		GetProfile(ctx context.Context, id uuid.UUID) (*v1u.GetProfileResponse, error)
	}
	PostDomainService interface {
		GetPostNormal(ctx context.Context, id uuid.UUID) (*v1p.GetPostNormalResponse, error)
	}
)

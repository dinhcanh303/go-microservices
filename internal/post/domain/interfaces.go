package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type (
	CommentDomainService interface {
		GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*sharedkernel.CommentHasChildren, error)
		CountCommentByPostID(ctx context.Context, postId uuid.UUID) (int64, error)
	}
	LikeDomainService interface {
		GetLikesByPostID(ctx context.Context, postId uuid.UUID) (*domain.LikesInfo, error)
	}
	UploadDomainService interface {
		GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domainUpload.Attachment, error)
	}
	GroupDomainService interface {
		GetAllGroupIdByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	}
	AuthDomainService interface {
		GetAllUserIdByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	}
)

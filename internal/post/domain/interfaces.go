package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	domain2 "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type (
	CommentDomainService interface {
		GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*sharedkernel.CommentHasChildren, error)
	}
	LikeDomainService interface {
		GetLikesByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Like, error)
	}
	UploadDomainService interface {
		GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domain2.Attachment, error)
	}
)

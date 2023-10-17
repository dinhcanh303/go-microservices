package domain

import (
	"context"

	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type (
	CommentDomainService interface {
		GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*sharedkernel.CommentHasChildren, error)
	}
	LikeDomainService interface {
		GetLikesByPostID(ctx context.Context, postId uuid.UUID) ([]*LikeItem, error)
	}
)

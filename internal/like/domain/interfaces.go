package domain

import (
	"context"

	v1c "github.com/dinhcanh303/go-microservices/api/comment/v1"
	v1p "github.com/dinhcanh303/go-microservices/api/post/v1"
	"github.com/google/uuid"
)

type (
	PostDomainService interface {
		GetPostNormal(ctx context.Context, id uuid.UUID) (*v1p.GetPostNormalResponse, error)
	}
	CommentDomainService interface {
		GetComment(ctx context.Context, id uuid.UUID) (*v1c.GetCommentResponse, error)
	}
)

package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
)

type (
	PostDomainService interface {
		GetPostNormal(ctx context.Context, id uuid.UUID) (*gen.GetPostNormalResponse, error)
	}
	CommentDomainService interface {
		GetComment(ctx context.Context, id uuid.UUID) (*gen.GetCommentResponse, error)
	}
)

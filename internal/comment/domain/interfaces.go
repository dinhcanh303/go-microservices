package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/google/uuid"
)

type LikeDomainService interface {
	GetLikesByCommentID(ctx context.Context, commentId uuid.UUID) ([]*domain.Like, error)
}

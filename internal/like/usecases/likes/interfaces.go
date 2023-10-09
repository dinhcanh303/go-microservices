package likes

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/google/uuid"
)

type LikeRepo interface {
	Create(ctx context.Context, like *domain.Like) (*domain.Like, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Like, error)
	Update(ctx context.Context, like *domain.Like) (*domain.Like, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	CountByPostID(ctx context.Context, postId uuid.UUID) (uint64, error)
	CountByCommentID(ctx context.Context, commentId uuid.UUID) (uint64, error)
}
type UseCase interface {
	CreateLike(ctx context.Context, like *domain.Like) (*domain.Like, error)
	GetLike(ctx context.Context, id uuid.UUID) (*domain.Like, error)
	UpdateLike(ctx context.Context, like *domain.Like) (*domain.Like, error)
	DeleteLike(ctx context.Context, id uuid.UUID) (bool, error)
	CountLikeByPostID(ctx context.Context, postId uuid.UUID) (uint64, error)
	CountLikeByCommentID(ctx context.Context, commentId uuid.UUID) (uint64, error)
}

package likes

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/google/uuid"
)

type LikeRepo interface {
	Create(ctx context.Context, like *domain.Like) (*domain.Like, error)
	Update(ctx context.Context, like *domain.Like) (*domain.Like, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	GetLikeByUserId(ctx context.Context, likeableType string, likeableId, userId uuid.UUID) (*domain.Like, error)
	GetLikesByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Like, error)
	GetLikesByCommentID(ctx context.Context, commentID uuid.UUID) ([]*domain.Like, error)
	GetLikesInfoByPostID(ctx context.Context, postID, userID uuid.UUID) (*domain.LikesInfo, error)
	GetLikesInfoByCommentID(ctx context.Context, commentID, userID uuid.UUID) (*domain.LikesInfo, error)
}
type UseCase interface {
	CreateLike(ctx context.Context, like *domain.Like) (*domain.Like, error)
	UpdateLike(ctx context.Context, like *domain.Like) (*domain.Like, error)
	DeleteLike(ctx context.Context, id uuid.UUID) (bool, error)
	GetLikeByUserId(ctx context.Context, likeableType string, likeableId, userId uuid.UUID) (*domain.Like, error)
	GetLikesByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Like, error)
	GetLikesByCommentID(ctx context.Context, commentID uuid.UUID) ([]*domain.Like, error)
	GetLikesInfoByPostID(ctx context.Context, postID, userID uuid.UUID) (*domain.LikesInfo, error)
	GetLikesInfoByCommentID(ctx context.Context, commentID, userID uuid.UUID) (*domain.LikesInfo, error)
}

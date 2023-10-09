package likes

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/google/uuid"
)

type service struct {
	likeRepo LikeRepo
}

// CountLikeByCommentID implements UseCase.
func (*service) CountLikeByCommentID(ctx context.Context, commentId uuid.UUID) (uint64, error) {
	panic("unimplemented")
}

// CountLikeByPostID implements UseCase.
func (*service) CountLikeByPostID(ctx context.Context, postId uuid.UUID) (uint64, error) {
	panic("unimplemented")
}

// CreateLike implements UseCase.
func (*service) CreateLike(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	panic("unimplemented")
}

// DeleteLike implements UseCase.
func (*service) DeleteLike(ctx context.Context, id uuid.UUID) (bool, error) {
	panic("unimplemented")
}

// GetLike implements UseCase.
func (*service) GetLike(ctx context.Context, id uuid.UUID) (*domain.Like, error) {
	panic("unimplemented")
}

// UpdateLike implements UseCase.
func (*service) UpdateLike(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	panic("unimplemented")
}

var _ UseCase = (*service)(nil)

func NewService(likeRepo LikeRepo) UseCase {
	return &service{
		likeRepo: likeRepo,
	}
}

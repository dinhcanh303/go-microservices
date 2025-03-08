package follow

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/google/uuid"
)

type (
	FollowRepo interface {
		CreateFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error
		DeleteFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error
		GetFollowers(ctx context.Context, followingId uuid.UUID) ([]*domain.UserFollow, error)
		GetFollowing(ctx context.Context, followerId uuid.UUID) ([]*domain.UserFollow, error)
		GetFollowingIds(ctx context.Context, followerId uuid.UUID) ([]uuid.UUID, error)
	}
	UseCase interface {
		Follow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error
		UnFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error
		GetFollowers(ctx context.Context, followingId uuid.UUID) ([]*domain.UserFollow, error)
		GetFollowing(ctx context.Context, followerId uuid.UUID) ([]*domain.UserFollow, error)
		GetFollowingIds(ctx context.Context, followerId uuid.UUID) ([]uuid.UUID, error)
	}
)

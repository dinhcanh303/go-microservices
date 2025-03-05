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
		GetFollowers(context.Context, uuid.UUID) ([]*domain.UserFollow, error)
		GetFollowing(context.Context, uuid.UUID) ([]*domain.UserFollow, error)
	}
	UseCase interface {
		Follow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error
		UnFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error
		GetFollowers(context.Context, uuid.UUID) ([]*domain.UserFollow, error)
		GetFollowing(context.Context, uuid.UUID) ([]*domain.UserFollow, error)
	}
)

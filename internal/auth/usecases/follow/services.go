package follow

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/google/uuid"
	"github.com/google/wire"
)

type service struct {
	repo  FollowRepo
	redis redis.RedisEngine
}

// Follow implements UseCase.
func (s *service) Follow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error {
	panic("unimplemented")
}

// GetFollowers implements UseCase.
func (s *service) GetFollowers(context.Context, uuid.UUID) ([]*domain.UserFollow, error) {
	panic("unimplemented")
}

// GetFollowing implements UseCase.
func (s *service) GetFollowing(context.Context, uuid.UUID) ([]*domain.UserFollow, error) {
	panic("unimplemented")
}

// UnFollow implements UseCase.
func (s *service) UnFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error {
	panic("unimplemented")
}

var _ UseCase = (*service)(nil)

func NewUseCase(repo FollowRepo, redis redis.RedisEngine) UseCase {
	return &service{
		repo:  repo,
		redis: redis,
	}
}

var UseCaseSet = wire.NewSet(NewUseCase)

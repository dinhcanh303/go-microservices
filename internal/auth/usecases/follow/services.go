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

// GetFollowingIds implements UseCase.
func (s *service) GetFollowingIds(ctx context.Context, followerId uuid.UUID) ([]uuid.UUID, error) {
	res, err := s.repo.GetFollowingIds(ctx, followerId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Follow implements UseCase.
func (s *service) Follow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error {
	err := s.repo.CreateFollow(ctx, followerId, followingId)
	if err != nil {
		return err
	}
	return nil
}

// GetFollowers implements UseCase.
func (s *service) GetFollowers(ctx context.Context, followingId uuid.UUID) ([]*domain.UserFollow, error) {
	followers, err := s.repo.GetFollowers(ctx, followingId)
	if err != nil {
		return nil, err
	}
	return followers, nil
}

// GetFollowing implements UseCase.
func (s *service) GetFollowing(ctx context.Context, followerId uuid.UUID) ([]*domain.UserFollow, error) {
	following, err := s.repo.GetFollowing(ctx, followerId)
	if err != nil {
		return nil, err
	}
	return following, nil
}

// UnFollow implements UseCase.
func (s *service) UnFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error {
	err := s.repo.DeleteFollow(ctx, followerId, followingId)
	if err != nil {
		return err
	}
	return nil
}

var _ UseCase = (*service)(nil)

func NewUseCase(repo FollowRepo, redis redis.RedisEngine) UseCase {
	return &service{
		repo:  repo,
		redis: redis,
	}
}

var UseCaseSet = wire.NewSet(NewUseCase)

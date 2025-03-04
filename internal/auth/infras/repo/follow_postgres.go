package repo

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/follow"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
)

type followRepo struct {
	pg postgres.DBEngine
}

// CreateFollow implements follow.FollowRepo.
func (f *followRepo) CreateFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) (*domain.Follow, error) {
	panic("unimplemented")
}

// DeleteFollow implements follow.FollowRepo.
func (f *followRepo) DeleteFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) (*domain.Follow, error) {
	panic("unimplemented")
}

// GetFollowers implements follow.FollowRepo.
func (f *followRepo) GetFollowers(context.Context, uuid.UUID) ([]*domain.Follow, error) {
	panic("unimplemented")
}

// GetFollowing implements follow.FollowRepo.
func (f *followRepo) GetFollowing(context.Context, uuid.UUID) ([]*domain.Follow, error) {
	panic("unimplemented")
}

var _ follow.FollowRepo = (*followRepo)(nil)

func NewFollowRepo(pg postgres.DBEngine) follow.FollowRepo {
	return &followRepo{
		pg: pg,
	}
}

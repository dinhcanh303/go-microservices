package repo

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/follow"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type followRepo struct {
	pg postgres.DBEngine
}

// CreateFollow implements follow.FollowRepo.
func (rp *followRepo) CreateFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "followRepo.CreateFollow db failed")
	}
	qtx := querier.WithTx(tx)
	_, err = qtx.CreateFollow(ctx, postgresql.CreateFollowParams{
		FollowerID:  followerId,
		FollowingID: followingId,
	})
	if err != nil {
		return errors.Wrap(err, "followRepo.CreateFollow failed")
	}
	return tx.Commit()
}

// DeleteFollow implements follow.FollowRepo.
func (rp *followRepo) DeleteFollow(ctx context.Context, followerId uuid.UUID, followingId uuid.UUID) error {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "followRepo.DeleteFollow db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.DeleteFollow(ctx, postgresql.DeleteFollowParams{
		FollowerID:  followerId,
		FollowingID: followingId,
	})
	if err != nil {
		return errors.Wrap(err, "followRepo.DeleteFollow failed")
	}
	return tx.Commit()
}

// GetFollowers implements follow.FollowRepo.
func (rp *followRepo) GetFollowers(ctx context.Context, followingId uuid.UUID) ([]*domain.UserFollow, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	res, err := querier.GetFollowers(ctx, followingId)
	if err != nil {
		return nil, errors.Wrap(err, "followRepo.GetFollowers failed")
	}
	return lo.Map(res, func(follower postgresql.GetFollowersRow, _ int) *domain.UserFollow {
		return &domain.UserFollow{
			Id:        follower.ID,
			NickName:  follower.NickName.String,
			FullName:  follower.FullName.String,
			AvatarUrl: follower.AvatarUrl.String,
		}
	}), nil
}

// GetFollowing implements follow.FollowRepo.
func (rp *followRepo) GetFollowing(ctx context.Context, followerId uuid.UUID) ([]*domain.UserFollow, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	res, err := querier.GetFollowing(ctx, followerId)
	if err != nil {
		return nil, errors.Wrap(err, "followRepo.GetFollowers failed")
	}
	return lo.Map(res, func(following postgresql.GetFollowingRow, _ int) *domain.UserFollow {
		return &domain.UserFollow{
			Id:        following.ID,
			NickName:  following.NickName.String,
			FullName:  following.FullName.String,
			AvatarUrl: following.AvatarUrl.String,
		}
	}), nil
}

var _ follow.FollowRepo = (*followRepo)(nil)

func NewFollowRepo(pg postgres.DBEngine) follow.FollowRepo {
	return &followRepo{
		pg: pg,
	}
}

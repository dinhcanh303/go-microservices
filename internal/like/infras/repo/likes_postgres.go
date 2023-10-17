package repo

import (
	"context"
	"database/sql"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/internal/like/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/like/usecases/likes"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type likeRepo struct {
	pg postgres.DBEngine
}

var _ likes.LikeRepo = (*likeRepo)(nil)

var RepositorySet = wire.NewSet(NewLikeRepo)

func NewLikeRepo(pg postgres.DBEngine) likes.LikeRepo {
	return &likeRepo{
		pg: pg,
	}
}

// Create implements likes.LikeRepo.
func (rp *likeRepo) Create(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "likeRepo.Create db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Create(ctx, postgresql.CreateParams{
		ID:           uuid.New(),
		Emoji:        like.Emoji,
		LikeableType: like.LikeableType,
		LikeableID:   like.LikeableID,
		UserID:       like.UserID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "likeRepo.Create failed")
	}
	return &domain.Like{
		ID:           result.ID,
		Emoji:        result.Emoji,
		LikeableType: result.LikeableType,
		LikeableID:   result.LikeableID,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}, tx.Commit()
}

// Delete implements likes.LikeRepo.
func (rp *likeRepo) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "likeRepo.Delete db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "likeRepo.Delete failed")
	}
	return true, tx.Commit()
}

// GetLikesByCommentID implements likes.LikeRepo.
func (rp *likeRepo) GetLikesByCommentID(ctx context.Context, commentID uuid.UUID) ([]*domain.Like, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	results, err := querier.GetAllByType(ctx, postgresql.GetAllByTypeParams{
		LikeableType: "Like/Comment",
		LikeableID:   commentID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "likeRepo.GetLikesByCommentID failed")
	}
	return lo.Map(results, func(item postgresql.LikeLike, _ int) *domain.Like {
		return &domain.Like{
			ID:           item.ID,
			Emoji:        item.Emoji,
			LikeableType: item.LikeableType,
			LikeableID:   item.LikeableID,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}), nil

}

// GetLikesByPostID implements likes.LikeRepo.
func (rp *likeRepo) GetLikesByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Like, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	results, err := querier.GetAllByType(ctx, postgresql.GetAllByTypeParams{
		LikeableType: "Like/Post",
		LikeableID:   postID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "likeRepo.GetLikesByPostID failed")
	}
	return lo.Map(results, func(item postgresql.LikeLike, _ int) *domain.Like {
		return &domain.Like{
			ID:           item.ID,
			Emoji:        item.Emoji,
			LikeableType: item.LikeableType,
			LikeableID:   item.LikeableID,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}), nil
}

// Update implements likes.LikeRepo.
func (rp *likeRepo) Update(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "likeRepo.Update db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Update(ctx, postgresql.UpdateParams{
		ID: like.ID,
		Emoji: sql.NullString{
			String: like.Emoji,
			Valid:  like.Emoji != "",
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "likeRepo.Update failed")
	}
	return &domain.Like{
		ID:           result.ID,
		Emoji:        result.Emoji,
		LikeableType: result.LikeableType,
		LikeableID:   result.LikeableID,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}, tx.Commit()
}

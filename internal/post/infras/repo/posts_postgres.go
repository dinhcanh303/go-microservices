package repo

import (
	"context"
	"database/sql"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type postRepo struct {
	pg postgres.DBEngine
}

func NewPostRepo(pg postgres.DBEngine) posts.PostRepo {
	return &postRepo{pg: pg}
}

var _ posts.PostRepo = (*postRepo)(nil)

var RepositoryPostSet = wire.NewSet(NewPostRepo)

// GetByGroupId implements posts.PostRepo.
func (rp *postRepo) GetByGroupId(ctx context.Context, groupId uuid.UUID) ([]*domain.Post, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	results, err := querier.GetByGroupId(ctx, uuid.NullUUID{
		UUID:  groupId,
		Valid: true, // Set this to true if groupId is valid, or false if it's NULL
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Get(ctx, id) failed")
	}
	return lo.Map(results, func(item postgresql.PostPost, _ int) *domain.Post {
		return &domain.Post{
			ID:        item.ID,
			Title:     item.Title,
			Content:   item.Content,
			Status:    item.Status,
			UserID:    item.UserID,
			GroupID:   item.GroupID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}), nil
}

// GetByUserId implements posts.PostRepo.
func (rp *postRepo) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Post, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	results, err := querier.GetByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Get(ctx, id) failed")
	}
	return lo.Map(results, func(item postgresql.PostPost, _ int) *domain.Post {
		return &domain.Post{
			ID:        item.ID,
			Title:     item.Title,
			Content:   item.Content,
			Status:    item.Status,
			UserID:    item.UserID,
			GroupID:   item.GroupID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}), nil
}

// List implements posts.PostRepo.
func (*postRepo) List(ctx context.Context, offset int, limit int) ([]*domain.Post, error) {
	panic("unimplemented")
}

// Create implements posts.PostRepo.
func (rp *postRepo) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Create db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Create(ctx, postgresql.CreateParams{
		ID:      uuid.New(),
		Title:   post.Title,
		Content: post.Content,
		Status:  post.Status,
		UserID:  post.UserID,
		GroupID: post.GroupID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Create(ctx, postgresql.CreateParams) failed")
	}

	return &domain.Post{
		ID:        result.ID,
		Title:     result.Title,
		Content:   result.Content,
		Status:    result.Status,
		UserID:    result.UserID,
		GroupID:   post.GroupID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, tx.Commit()
}

// Delete implements posts.PostRepo.
func (rp *postRepo) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "DeletePostRepo")
	}
	qtx := querier.WithTx(tx)
	err = qtx.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "qtx.Delete(ctx, id) failed")
	}
	return true, tx.Commit()
}

// Get implements posts.PostRepo.
func (rp *postRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	result, err := querier.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.Get(ctx, id) failed")
	}

	return &domain.Post{
		ID:        result.ID,
		Title:     result.Title,
		Content:   result.Content,
		Status:    result.Status,
		UserID:    result.UserID,
		GroupID:   result.GroupID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

// Update implements posts.PostRepo.
func (rp *postRepo) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "CreatePostRepo")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Update(ctx, postgresql.UpdateParams{
		ID: post.ID,
		Title: sql.NullString{
			String: post.Title,
			Valid:  post.Title != "",
		},
		Content: sql.NullString{
			String: post.Content,
			Valid:  post.Content != "",
		},
		Status: sql.NullInt32{
			Int32: post.Status,
			Valid: post.Status != 0,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Create(ctx, postgresql.UpdateParams) failed")
	}

	return &domain.Post{
		ID:        result.ID,
		Title:     result.Title,
		Content:   result.Content,
		Status:    result.Status,
		UserID:    result.UserID,
		GroupID:   result.GroupID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, tx.Commit()
}

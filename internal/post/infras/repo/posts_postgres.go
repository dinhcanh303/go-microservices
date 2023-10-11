package repo

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
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
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "GetPostRepo")
	}
	// Convert groupId to a NullUUID
	nullGroupId := uuid.NullUUID{
		UUID:  groupId,
		Valid: true, // Set this to true if groupId is valid, or false if it's NULL
	}
	qtx := querier.WithTx(tx)
	results, err := qtx.GetByGroupId(ctx, nullGroupId)
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
			DeletedAt: item.DeletedAt.Time,
		}
	}), tx.Commit()
}

// GetByUserId implements posts.PostRepo.
func (rp *postRepo) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Post, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "GetPostRepo")
	}
	qtx := querier.WithTx(tx)
	results, err := qtx.GetByUserId(ctx, userId)
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
			DeletedAt: item.DeletedAt.Time,
		}
	}), tx.Commit()
}

// List implements posts.PostRepo.
func (*postRepo) List(ctx context.Context, offset int, limit int) ([]*domain.Post, error) {
	panic("unimplemented")
}

// Create implements posts.PostRepo.
func (rp *postRepo) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	slog.Info("Repo Postgresql")
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "CreatePostRepo")
	}
	qtx := querier.WithTx(tx)
	slog.Info("QTX")
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
		DeletedAt: result.DeletedAt.Time,
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
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "GetPostRepo")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Get(ctx, id) failed")
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
		DeletedAt: result.DeletedAt.Time,
	}, tx.Commit()
}

// // GetWithUnscoped implements posts.PostRepo.
// func (rp *postRepo) GetWithUnscoped(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
// 	db := rp.pg.GetDB()
// 	querier := postgresql.New(db)
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "GetPostRepo")
// 	}
// 	qtx := querier.WithTx(tx)
// 	result, err := qtx.GetWithUnscoped(ctx, id)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "qtx.GetWithUnscoped(ctx, id) failed")
// 	}
// 	return &domain.Post{
// 		ID:          result.ID,
// 		Name:        result.Name,
// 		Description: result.Description,
// 		Status:      result.Status,
// 		UserID:      result.UserID,
// 		CreatedAt:   result.CreatedAt,
// 		UpdatedAt:   result.UpdatedAt,
// 		DeletedAt:   result.DeletedAt.Time,
// 	}, tx.Commit()
// }

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
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Status:  post.Status,
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
		GroupID:   result.GroupID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		DeletedAt: result.DeletedAt.Time,
	}, tx.Commit()
}

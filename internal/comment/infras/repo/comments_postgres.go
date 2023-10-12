package repo

import (
	"context"
	"database/sql"

	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/dinhcanh303/go-microservices/internal/comment/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type commentRepo struct {
	pg postgres.DBEngine
}

var _ comments.CommentRepo = (*commentRepo)(nil)

var RepositorySet = wire.NewSet(NewCommentRepo)

func NewCommentRepo(pg postgres.DBEngine) comments.CommentRepo {
	return &commentRepo{
		pg: pg,
	}
}

// CountByCommentID implements comments.CommentRepo.
func (rp *commentRepo) CountByCommentID(ctx context.Context, commentId uuid.UUID) (int64, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	result, err := querier.CountByCommentID(ctx, uuid.NullUUID{
		UUID:  commentId,
		Valid: true,
	})
	if err != nil {
		return 0, errors.Wrap(err, "commentRepo.CountByCommentID failed")
	}
	return result, nil
}

// CountByPostID implements comments.CommentRepo.
func (rp *commentRepo) CountByPostID(ctx context.Context, postId uuid.UUID) (int64, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	result, err := querier.CountByPostID(ctx, postId)
	if err != nil {
		return 0, errors.Wrap(err, "commentRepo.CountByPostID failed")
	}
	return result, nil
}

// Create implements comments.CommentRepo.
func (rp *commentRepo) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.Create db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Create(ctx, postgresql.CreateParams{
		ID:              comment.ID,
		UserID:          comment.UserID,
		Content:         comment.Content,
		PostID:          comment.PostID,
		ParentCommentID: comment.ParentCommentID,
		ReplyTo: sql.NullString{
			String: comment.ReplyTo,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.Create failed")
	}
	return &domain.Comment{
		ID:              result.ID,
		UserID:          result.UserID,
		Content:         result.Content,
		PostID:          result.PostID,
		ParentCommentID: result.ParentCommentID,
		ReplyTo:         result.ReplyTo.String,
		CreatedAt:       result.CreatedAt,
		UpdatedAt:       result.UpdatedAt,
	}, tx.Commit()
}

// Delete implements comments.CommentRepo.
func (rp *commentRepo) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "commentRepo.Delete db failed")
	}
	qtx := querier.WithTx(tx)

	err = qtx.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "commentRepo.Delete failed")
	}
	return true, tx.Commit()
}

// DeleteAllByPostID implements comments.CommentRepo.
func (rp *commentRepo) DeleteAllByPostID(ctx context.Context, postId uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "commentRepo.DeleteAllByPostID db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.DeleteAllByPostID(ctx, postId)
	if err != nil {
		return false, errors.Wrap(err, "commentRepo.DeleteAllByPostID failed")
	}
	return true, tx.Commit()
}

// Get implements comments.CommentRepo.
func (rp *commentRepo) Get(ctx context.Context, uuid uuid.UUID) (*domain.Comment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	result, err := querier.Get(ctx, uuid)
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.Get failed")
	}
	return &domain.Comment{
		ID:              result.ID,
		UserID:          result.UserID,
		Content:         result.Content,
		PostID:          result.PostID,
		ParentCommentID: result.ParentCommentID,
		ReplyTo:         result.ReplyTo.String,
		CreatedAt:       result.CreatedAt,
		UpdatedAt:       result.UpdatedAt,
	}, nil
}

// GetCommentByPostID implements comments.CommentRepo.
func (rp *commentRepo) GetCommentByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Comment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	results, err := querier.GetCommentByPostID(ctx, postId)
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.GetCommentByPostID failed")
	}
	return lo.Map(results, func(item postgresql.CommentComment, _ int) *domain.Comment {
		return &domain.Comment{
			ID:              item.ID,
			UserID:          item.UserID,
			Content:         item.Content,
			PostID:          item.PostID,
			ParentCommentID: item.ParentCommentID,
			ReplyTo:         item.ReplyTo.String,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
		}
	}), nil
}

// Update implements comments.CommentRepo.
func (rp *commentRepo) Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.Create db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Update(ctx, postgresql.UpdateParams{
		ID:      comment.ID,
		Content: comment.Content,
		ReplyTo: sql.NullString{
			String: comment.ReplyTo,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.Create failed")
	}
	return &domain.Comment{
		ID:              result.ID,
		UserID:          result.UserID,
		Content:         result.Content,
		PostID:          result.PostID,
		ParentCommentID: result.ParentCommentID,
		ReplyTo:         result.ReplyTo.String,
		CreatedAt:       result.CreatedAt,
		UpdatedAt:       result.UpdatedAt,
	}, tx.Commit()
}

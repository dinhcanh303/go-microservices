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
func (rp *commentRepo) CountByCommentID(ctx context.Context, commentId uuid.UUID) (uint64, error) {
	panic("unimplemented")
}

// CountByPostID implements comments.CommentRepo.
func (rp *commentRepo) CountByPostID(ctx context.Context, postId uuid.UUID) (uint64, error) {
	panic("unimplemented")
}

// Create implements comments.CommentRepo.
func (rp *commentRepo) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.Create")
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
		return nil, errors.Wrap(err, "commentRepo.Create")
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
		DeletedAt:       result.DeletedAt.Time,
	}, tx.Commit()
}

// Delete implements comments.CommentRepo.
func (rp *commentRepo) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "commentRepo.Create")
	}
	qtx := querier.WithTx(tx)

	err = qtx.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "commentRepo.Create")
	}
	return true, tx.Commit()
}

// DeleteByCommentID implements comments.CommentRepo.
func (rp *commentRepo) DeleteByCommentID(ctx context.Context, commentId uuid.UUID) (bool, error) {
	panic("unimplemented")
}

// Get implements comments.CommentRepo.
func (rp *commentRepo) Get(ctx context.Context, uuid uuid.UUID) (*domain.Comment, error) {
	panic("unimplemented")
}

// ListByPostID implements comments.CommentRepo.
func (rp *commentRepo) ListByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Comment, error) {
	panic("unimplemented")
}

// Update implements comments.CommentRepo.
func (rp *commentRepo) Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	panic("unimplemented")
}

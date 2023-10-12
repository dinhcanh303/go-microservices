package comments

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/google/uuid"
)

type (
	CommentRepo interface {
		Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
		Get(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
		Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
		Delete(ctx context.Context, id uuid.UUID) (bool, error)
		DeleteAllByPostID(ctx context.Context, postId uuid.UUID) (bool, error)
		GetCommentByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Comment, error)
		CountByPostID(ctx context.Context, postId uuid.UUID) (int64, error)
		CountByCommentID(ctx context.Context, commentId uuid.UUID) (int64, error)
	}
	UseCase interface {
		CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
		GetComment(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
		UpdateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
		DeleteComment(ctx context.Context, id uuid.UUID) (bool, error)
		DeleteAllCommentByPostID(ctx context.Context, postId uuid.UUID) (bool, error)
		GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Comment, error)
		CountCommentByPostID(ctx context.Context, postId uuid.UUID) (int64, error)
		CountCommentByCommentID(ctx context.Context, commentId uuid.UUID) (int64, error)
	}
)

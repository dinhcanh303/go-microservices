package comments

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/google/uuid"
)

type CommentRepo interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	DeleteByCommentID(ctx context.Context, commentId uuid.UUID) (bool, error)
	ListByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Comment, error)
	CountByPostID(ctx context.Context, postId uuid.UUID) (uint64, error)
	CountByCommentID(ctx context.Context, commentId uuid.UUID) (uint64, error)
}
type UseCase interface {
	CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetComment(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
	UpdateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) (bool, error)
	DeleteCommentByCommentID(ctx context.Context, commentId uuid.UUID) (bool, error)
	ListCommentByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Comment, error)
	CountCommentByPostID(ctx context.Context, postId uuid.UUID) (uint64, error)
	CountCommentByCommentID(ctx context.Context, commentId uuid.UUID) (uint64, error)
}

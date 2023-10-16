package posts

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	publisher "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/google/uuid"
)

type (
	PostRepo interface {
		Get(ctx context.Context, id uuid.UUID) (*domain.Post, error)
		GetByGroupId(ctx context.Context, groupId uuid.UUID) ([]*domain.Post, error)
		GetByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Post, error)
		Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
		Update(ctx context.Context, post *domain.Post) (*domain.Post, error)
		Delete(ctx context.Context, id uuid.UUID) (bool, error)
		List(ctx context.Context, offset, limit int) ([]*domain.Post, error)
	}
	UploadEventPublisher interface {
		Configure(...publisher.Option)
		Publish(context.Context, []byte, string) error
	}
	UseCase interface {
		GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error)
		GetPostExtra(ctx context.Context, id uuid.UUID) (*domain.PostExtra, error)
		GetByGroupId(ctx context.Context, groupId uuid.UUID) ([]*domain.Post, error)
		GetByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Post, error)
		CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
		UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
		DeletePost(ctx context.Context, id uuid.UUID) (bool, error)
		ListPost(ctx context.Context, offset, limit int) ([]*domain.Post, error)
	}
)

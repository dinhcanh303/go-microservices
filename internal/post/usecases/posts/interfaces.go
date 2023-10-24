package posts

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
)

type (
	PostRepo interface {
		Get(ctx context.Context, id uuid.UUID) (*domain.Post, error)
		Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
		Update(ctx context.Context, post *domain.Post) (*domain.Post, error)
		Delete(ctx context.Context, id uuid.UUID) (bool, error)
		GetByUserId(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]*domain.Post, error)
		GetByGroupId(ctx context.Context, groupId uuid.UUID, limit, offset int32) ([]*domain.Post, error)
		GetByFeed(ctx context.Context, userIds, groupIds string, limit, offset int32) ([]*domain.Post, error)
	}
	UseCase interface {
		GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error)
		CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
		UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
		DeletePost(ctx context.Context, id uuid.UUID) (bool, error)
		GetPostsByUserId(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]*domain.PostExtra, error)
		GetPostsByGroupId(ctx context.Context, groupId uuid.UUID, limit, offset int32) ([]*domain.PostExtra, error)
		GetPostsByFeed(ctx context.Context, userIds, groupIds string, limit, offset int32) ([]*domain.PostExtra, error)
	}
)

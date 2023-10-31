package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
)

type (
	PostDomainService interface {
		GetPostsByUserId(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]*domain.Post, error)
		GetPostsByGroupId(ctx context.Context, groupId uuid.UUID, limit, offset int32) ([]*domain.Post, error)
		GetPostsByFeed(ctx context.Context, userIds, groupIds string, limit, offset int32) ([]*domain.Post, error)
	}
)

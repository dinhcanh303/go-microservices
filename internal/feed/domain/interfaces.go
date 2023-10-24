package domain

import (
	"context"

	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type (
	PostDomainService interface {
		GetPostsByFeed(ctx context.Context, userIds, groupIds string, limit, offset int32) ([]*sharedkernel.CommentHasChildren, error)
	}
	GroupDomainService interface {
		GetAllGroupIdByUserId(ctx context.Context, userId uuid.UUID) ([]string, error)
	}
)

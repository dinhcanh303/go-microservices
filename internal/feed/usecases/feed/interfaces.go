package posts

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
)

type (
	UseCase interface {
		NewFeed(ctx context.Context, userId uuid.UUID) ([]*domain.PostExtra, error)
	}
)

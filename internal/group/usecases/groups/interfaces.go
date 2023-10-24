package groups

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/google/uuid"
)

type GroupRepo interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Group, error)
	Create(ctx context.Context, group *domain.Group) (*domain.Group, error)
	Update(ctx context.Context, group *domain.Group) (*domain.Group, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	GetAllGroupByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Group, error)
	GetAllGroupIdByUserId(ctx context.Context, userId uuid.UUID) ([]string, error)
}

type UseCase interface {
	GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, error)
	CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error)
	UpdateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error)
	DeleteGroup(ctx context.Context, id uuid.UUID) (bool, error)
	GetAllGroupByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Group, error)
	GetAllGroupIdByUserId(ctx context.Context, userId uuid.UUID) ([]string, error)
}

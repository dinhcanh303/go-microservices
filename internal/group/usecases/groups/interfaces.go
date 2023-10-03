package groups

import (
	"context"
	"go-microservices/internal/group/domain"
)

type UseCase interface {
	GetGroup(ctx context.Context, uuid string) (*domain.Group, error)
	GetGroupWithUnscoped(ctx context.Context, uuid string) (*domain.Group, error)
	CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error)
	UpdateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error)
	DeleteGroup(ctx context.Context, uuid string) (bool, error)
}

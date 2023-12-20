package notifications

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	"github.com/google/uuid"
)

type NotificationRepo interface {
	GetAllByID(context.Context, uuid.UUID) ([]*domain.Notification, error)
	GetListByID(ctx context.Context, id uuid.UUID, offset, limit uint32) ([]*domain.Notification, error)
	Update(context.Context, *domain.Notification) error
}

type UseCase interface {
	GetListNotificationByID(ctx context.Context, id uuid.UUID, offset, limit uint32) ([]*domain.Notification, error)
}

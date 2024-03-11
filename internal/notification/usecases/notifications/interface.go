package notifications

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type NotificationRepo interface {
	GetNotificationsByUserId(ctx context.Context, userId uuid.UUID, options sharedkernel.GetNotiOptions) ([]domain.Notification, error)
	CreateNotification(ctx context.Context, noti *domain.Notification) error
	UpdateNotification(ctx context.Context, id int, noti *domain.Notification) error
}

type UseCase interface {
	GetNotificationsByUserId(ctx context.Context, userId uuid.UUID, options sharedkernel.GetNotiOptions) error
	ReadNotification(ctx context.Context, noti *domain.Notification) error
}

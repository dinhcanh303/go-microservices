package notifications

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	"github.com/google/uuid"
)

type NotificationRepo interface {
	GetNotiByUserId(ctx context.Context, userId uuid.UUID, options string) ([]*domain.Notification, error)
	UpsertNoti(ctx context.Context, noti *domain.Notification) error
}

type UseCase interface {
	UpsertNoti(ctx context.Context, noti *domain.Notification) error
	GetNotiReadByUserId(ctx context.Context, userId uuid.UUID) error
	GetNotiUnReadByUserId(ctx context.Context, userId uuid.UUID) error
}

package notifications

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	repo  NotificationRepo
	redis redis.RedisEngine
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(
	repo NotificationRepo,
	redis redis.RedisEngine,
) UseCase {
	return &service{
		repo:  repo,
		redis: redis,
	}
}

// GetNotificationsByUserId implements UseCase.
func (s *service) GetNotificationsByUserId(ctx context.Context, userId uuid.UUID, limit, offset int, options sharedkernel.GetNotiOptions) ([]domain.Notification, error) {
	var notifications []domain.Notification
	keyCache := constant.CacheNotifications + userId.String() + constant.CacheLimit +
		utils.String(int32(limit)) + constant.CacheOffset + utils.String(int32(offset))
	err := utils.HandleHitCache(notifications, s.redis, keyCache)
	if err != nil {
		notifications, err = s.repo.GetNotificationsByUserId(ctx, userId, limit, offset, options)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetNotificationsByUserId failed")
		}
		err = s.redis.Set(keyCache, notifications)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return notifications, nil
}

// ReadNotification implements UseCase.
func (s *service) ReadNotification(ctx context.Context, noti *domain.Notification) error {
	err := s.repo.UpdateNotification(ctx, noti)
	if err != nil {
		return errors.Wrap(err, "service.CreateNotification")
	}
	err = s.redis.InvalidatePrefix(constant.CacheNotifications + noti.SenderID)
	if err != nil {
		slog.Error("InvalidatePrefix cache key failed")
	}
	return nil
}

// CreateNotification implements UseCase.
func (s *service) CreateNotification(ctx context.Context, noti *domain.Notification) error {
	err := s.repo.CreateNotification(ctx, noti)
	if err != nil {
		return errors.Wrap(err, "service.CreateNotification")
	}
	err = s.redis.InvalidatePrefix(constant.CacheNotifications + noti.SenderID)
	if err != nil {
		slog.Error("InvalidatePrefix cache key failed")
	}
	return nil
}

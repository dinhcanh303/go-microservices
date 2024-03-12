package notifications

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
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
	results, err := s.repo.GetNotificationsByUserId(ctx, userId, limit, offset, options)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateNotification")
	}
	return results, nil
}

// ReadNotification implements UseCase.
func (s *service) ReadNotification(ctx context.Context, noti *domain.Notification) error {
	err := s.repo.UpdateNotification(ctx, noti)
	if err != nil {
		return errors.Wrap(err, "service.CreateNotification")
	}
	return nil
}

// CreateNotification implements UseCase.
func (s *service) CreateNotification(ctx context.Context, noti *domain.Notification) error {
	err := s.repo.CreateNotification(ctx, noti)
	if err != nil {
		return errors.Wrap(err, "service.CreateNotification")
	}
	return nil
}

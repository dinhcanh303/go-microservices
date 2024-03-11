package notifications

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
	"github.com/google/wire"
)

type service struct {
	repo NotificationRepo
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(repo NotificationRepo) UseCase {
	return &service{
		repo: repo,
	}
}

// GetNotificationsByUserId implements UseCase.
func (*service) GetNotificationsByUserId(ctx context.Context, userId uuid.UUID, options sharedkernel.GetNotiOptions) error {
	panic("unimplemented")
}

// ReadNotification implements UseCase.
func (*service) ReadNotification(ctx context.Context, noti *domain.Notification) error {
	panic("unimplemented")
}

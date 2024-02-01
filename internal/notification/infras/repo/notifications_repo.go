package repo

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	"github.com/dinhcanh303/go-microservices/internal/notification/usecases/notifications"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type notificationRepo struct {
	db *gorm.DB
}

// GetNotiByUserId implements notifications.NotificationRepo.
func (rp *notificationRepo) GetNotiByUserId(ctx context.Context, userId uuid.UUID, options string) ([]*domain.Notification, error) {
	var results []*domain.Notification
	db := rp.db.Table("noti.notifications").Where(
		"receiver_id @> ? AND NOT ? <@read_id",
		[]string{userId.String()}, []string{userId.String()})
	if err := db.Select("*").Offset(0).Limit(20).Find(&results).Error; err != nil {
		return nil, errors.Wrap(err, "repo.GetNotiByUserId failed")
	}
	return results, nil

}

// UpsertNoti implements notifications.NotificationRepo.
func (rp *notificationRepo) UpsertNoti(ctx context.Context, noti *domain.Notification) error {
	if err := rp.db.Table("noti.notifications").Create(noti).Error; err != nil {
		return errors.Wrap(err, "create noti failed")
	}
	return nil
}

var RepositorySet = wire.NewSet()

func NewNotificationRepo(db *gorm.DB) notifications.NotificationRepo {
	return &notificationRepo{
		db: db,
	}
}

var _ notifications.NotificationRepo = (*notificationRepo)(nil)

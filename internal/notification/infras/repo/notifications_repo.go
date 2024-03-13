package repo

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	"github.com/dinhcanh303/go-microservices/internal/notification/usecases/notifications"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type notificationRepo struct {
	db *gorm.DB
}

var _ notifications.NotificationRepo = (*notificationRepo)(nil)

var RepositoryNotiSet = wire.NewSet(NewNotificationRepo)

func NewNotificationRepo(db *gorm.DB) notifications.NotificationRepo {
	return &notificationRepo{
		db: db,
	}
}

// GetNotiByUserId implements notifications.NotificationRepo.
func (rp *notificationRepo) GetNotificationsByUserId(ctx context.Context, userId uuid.UUID, limit, offset int, options sharedkernel.GetNotiOptions) ([]domain.Notification, error) {
	var results []domain.Notification
	db := rp.db.Table(domain.Notification{}.TableName()).Where(
		"sender_id =?", userId.String())
	if options.Unread {
		db = db.Where("read_at IS NULL")
	} else {
		db = db.Where("read_at IS NOT NULL")
	}
	if err := db.Select("*").Offset(offset).Limit(limit).Order("created_at DESC").Find(&results).Error; err != nil {
		return nil, errors.Wrap(err, "repo.GetNotiByUserId failed")
	}
	return results, nil

}

// UpsertNoti implements notifications.NotificationRepo.
func (rp *notificationRepo) CreateNotification(ctx context.Context, noti *domain.Notification) error {
	if err := rp.db.Table(noti.TableName()).Create(noti).Error; err != nil {
		return errors.Wrap(err, "create noti failed")
	}
	return nil
}

func (rp *notificationRepo) UpdateNotification(ctx context.Context, noti *domain.Notification) error {
	if err := rp.db.Table(noti.TableName()).Where("id = ?", noti.ID).Updates(noti).Error; err != nil {
		return errors.Wrap(err, "update noti failed")
	}
	return nil
}

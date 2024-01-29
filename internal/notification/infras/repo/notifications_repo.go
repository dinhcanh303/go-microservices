package repo

import (
	"github.com/dinhcanh303/go-microservices/internal/notification/usecases/notifications"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/wire"
)

type notificationRepo struct {
	pg postgres.DBEngine
}

var RepositorySet = wire.NewSet()

func NewNotificationRepo(pg postgres.DBEngine) notifications.NotificationRepo {
	return &notificationRepo{
		pg: pg,
	}
}

var _ notifications.NotificationRepo = (*notificationRepo)(nil)

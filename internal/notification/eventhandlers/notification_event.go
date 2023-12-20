package eventhandlers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
)

var _ NotificationEventHandler = (*notificationEventHandler)(nil)

type notificationEventHandler struct {
}

// Handle implements NotificationEventHandler.
func (*notificationEventHandler) Handle(context.Context, event.Notification) error {
	panic("unimplemented")
}

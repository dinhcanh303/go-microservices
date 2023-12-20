package eventhandlers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
)

type NotificationEventHandler interface {
	Handle(context.Context, event.Notification) error
}

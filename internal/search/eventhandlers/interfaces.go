package eventhandlers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
)

type EventHandlers interface {
	HandleUserCreated(context.Context, event.UserCreated) error
	HandleUserDeleted(context.Context, event.UserDeleted) error
	HandleGroupCreated(context.Context, event.GroupCreated) error
	HandleGroupDeleted(context.Context, event.GroupDeleted) error
}

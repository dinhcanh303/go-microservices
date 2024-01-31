package eventhandlers

import (
	"context"
)

type EventHandlers interface {
	HandleChangeDBUser(context.Context) error
	HandleChangeDBGroup(context.Context) error
}

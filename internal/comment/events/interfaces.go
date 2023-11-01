package events

import "context"

type (
	UserNotificationEventHandler interface {
		Handle(context.Context,*event.)
	}
)
package eventhandlers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/pkg/events"
)

type NotificationEventHandler interface {
	HandlerCommentNoti(context.Context, *events.Noti) error
	HandlerLikeNoti(context.Context, *events.Noti) error
	HandlerPostNoti(context.Context, *events.Noti) error
}

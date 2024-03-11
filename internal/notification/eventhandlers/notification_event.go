package eventhandlers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
)

var _ NotificationEventHandler = (*notificationEventHandler)(nil)

type notificationEventHandler struct {
}

// HandlerCommentNoti implements NotificationEventHandler.
func (*notificationEventHandler) HandlerCommentNoti(context.Context, *event.CommentNoti) error {
	panic("unimplemented")
}

// HandlerLikeNoti implements NotificationEventHandler.
func (*notificationEventHandler) HandlerLikeNoti(context.Context, *event.LikeNoti) error {
	panic("unimplemented")
}

// HandlerPostNoti implements NotificationEventHandler.
func (*notificationEventHandler) HandlerPostNoti(context.Context, *event.PostNoti) error {
	panic("unimplemented")
}

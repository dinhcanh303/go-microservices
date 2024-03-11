package eventhandlers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
)

type NotificationEventHandler interface {
	HandlerCommentNoti(context.Context, *event.CommentNoti) error
	HandlerLikeNoti(context.Context, *event.LikeNoti) error
	HandlerPostNoti(context.Context, *event.PostNoti) error
}

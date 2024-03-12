package eventhandlers

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	"github.com/dinhcanh303/go-microservices/internal/notification/usecases/notifications"
	"github.com/dinhcanh303/go-microservices/internal/pkg/events"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type eventhandlers struct {
	uc notifications.UseCase
}

var _ NotificationEventHandler = (*eventhandlers)(nil)

var EventHandlersSet = wire.NewSet(NewEventHandlers)

func NewEventHandlers(
	uc notifications.UseCase,
) NotificationEventHandler {
	return &eventhandlers{
		uc: uc,
	}
}

// HandlerCommentNoti implements NotificationEventHandler.
func (e *eventhandlers) HandlerCommentNoti(ctx context.Context, event *events.Noti) error {
	return createNotification(ctx, event, e.uc)
}

// HandlerLikeNoti implements NotificationEventHandler.
func (e *eventhandlers) HandlerLikeNoti(ctx context.Context, event *events.Noti) error {
	return createNotification(ctx, event, e.uc)
}

// HandlerPostNoti implements NotificationEventHandler.
func (e *eventhandlers) HandlerPostNoti(ctx context.Context, event *events.Noti) error {
	return createNotification(ctx, event, e.uc)
}

func createNotification(ctx context.Context, event *events.Noti, uc notifications.UseCase) error {
	checkErr := false
	var err error
	for _, senderId := range event.SenderIDs {
		err = uc.CreateNotification(ctx, &domain.Notification{
			ActorID:    event.ActorID,
			SenderID:   senderId,
			Data:       event.Data,
			Type:       event.Type,
			ObjectType: event.ObjectType,
			ObjectID:   event.ObjectID,
			ReadAt:     nil,
		})
		if err != nil {
			checkErr = true
			slog.Error("create notification failed:", err)
		}
	}
	if checkErr {
		return errors.Wrap(err, "Create Notification handler failed")
	}
	return nil
}

package app

import (
	"context"
	"encoding/json"

	"github.com/dinhcanh303/go-microservices/cmd/notification/config"
	"github.com/dinhcanh303/go-microservices/internal/notification/eventhandlers"
	"github.com/dinhcanh303/go-microservices/internal/pkg/events"
	consumer "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/consumer"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type App struct {
	Cfg            *config.Config
	PG             *gorm.DB
	NotiGRPCServer gen.NotiServiceServer
	AmqpConn       *amqp091.Connection
	Handlers       eventhandlers.NotificationEventHandler
	Consumer       consumer.EventConsumer
}

func New(
	cfg *config.Config,
	pg *gorm.DB,
	amqpConn *amqp091.Connection,
	notiGRPCServer gen.NotiServiceServer,
	consumer consumer.EventConsumer,
	handlers eventhandlers.NotificationEventHandler,
) *App {
	return &App{
		Cfg:            cfg,
		PG:             pg,
		AmqpConn:       amqpConn,
		Consumer:       consumer,
		Handlers:       handlers,
		NotiGRPCServer: notiGRPCServer,
	}
}

func (n *App) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)
		switch delivery.Type {
		case "post-noti":
			var payload events.Noti
			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}
			err = n.Handlers.HandlerPostNoti(ctx, &payload)
			checkLogErr(err, delivery)
		case "comment-noti":
			var payload events.Noti
			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}
			err = n.Handlers.HandlerCommentNoti(ctx, &payload)
			checkLogErr(err, delivery)
		case "like-noti":
			var payload events.Noti
			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}
			err = n.Handlers.HandlerLikeNoti(ctx, &payload)
			checkLogErr(err, delivery)
		}
	}
}
func checkLogErr(err error, delivery amqp091.Delivery) {
	if err != nil {
		if err = delivery.Reject(false); err != nil {
			slog.Error("failed to delivery.Reject", err)
		}

		slog.Error("failed to process delivery", err)
	} else {
		err = delivery.Ack(false)
		if err != nil {
			slog.Error("failed to acknowledge delivery", err)
		}
	}
}

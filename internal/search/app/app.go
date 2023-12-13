package app

import (
	"context"
	"encoding/json"

	"github.com/dinhcanh303/go-microservices/cmd/search/config"
	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	"github.com/dinhcanh303/go-microservices/internal/search/eventhandlers"
	consumer "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/comsumer"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

type App struct {
	Cfg              *config.Config
	AmqpConn         *amqp091.Connection
	Consumer         consumer.EventConsumer
	handlers         eventhandlers.EventHandlers
	SearchGRPCServer gen.SearchServiceServer
}

func New(
	cfg *config.Config,
	amqpConn *amqp091.Connection,
	consumer consumer.EventConsumer,
	handlers eventhandlers.EventHandlers,
	searchGRPCServer gen.SearchServiceServer,
) *App {
	return &App{
		Cfg:              cfg,
		AmqpConn:         amqpConn,
		Consumer:         consumer,
		handlers:         handlers,
		SearchGRPCServer: searchGRPCServer,
	}
}

func (c *App) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)
		switch delivery.Type {
		case "user-created":
			var payload event.UserCreated
			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal", err)
			}
			err = c.handlers.HandleUserCreated(ctx, payload)
			checkLogErr(err, delivery)
		case "user-deleted":
			var payload event.UserDeleted
			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal", err)
			}
			err = c.handlers.HandleUserDeleted(ctx, payload)
			checkLogErr(err, delivery)
		case "group-created":
			var payload event.GroupCreated
			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal", err)
			}
			err = c.handlers.HandleGroupCreated(ctx, payload)
			checkLogErr(err, delivery)
		case "group-deleted":
			var payload event.GroupDeleted
			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal", err)
			}
			payload.Identity()
			err = c.handlers.HandleGroupDeleted(ctx, payload)
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

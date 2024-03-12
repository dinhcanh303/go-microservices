package app

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/search/config"
	"github.com/dinhcanh303/go-microservices/internal/search/eventhandlers"
	consumer "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/consumer"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

type App struct {
	Cfg      *config.Config
	AmqpConn *amqp091.Connection
	Consumer consumer.EventConsumer
	// KafkaConsumer    kafka.KafkaConsumer
	handlers         eventhandlers.EventHandlers
	SearchGRPCServer gen.SearchServiceServer
}

func New(
	cfg *config.Config,
	amqpConn *amqp091.Connection,
	consumer consumer.EventConsumer,
	handlers eventhandlers.EventHandlers,
	searchGRPCServer gen.SearchServiceServer,
	// kafkaConsumer kafka.KafkaConsumer,
) *App {
	return &App{
		Cfg:              cfg,
		AmqpConn:         amqpConn,
		Consumer:         consumer,
		handlers:         handlers,
		SearchGRPCServer: searchGRPCServer,
		// KafkaConsumer:    kafkaConsumer,
	}
}

func (a *App) WorkerKafka() {
	// kafka.CheckConnector(kafka.ConnectorConfig{
	// 	ServerUrl:       "http://localhost:8083",
	// 	DataType:        "application/json",
	// 	ConnectorConfig: "./debezium/connectors/group_service_connector.json",
	// 	Connector:       "group_service_connector",
	// })
	// err := a.KafkaConsumer.Subscribe("postgres.group.groups", nil)
	// if err != nil {
	// 	slog.Error("Subscribe failed: ", err)
	// }
	// a.KafkaConsumer.ReadMessage(-1)
}

func (c *App) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)
		switch delivery.Type {
		case "users-changed":
			err := c.handlers.HandleChangeDBUser(ctx)
			checkLogErr(err, delivery)
		case "groups-changed":
			err := c.handlers.HandleChangeDBGroup(ctx)
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

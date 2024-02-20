package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	consumer "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/consumer"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestConnectRabbitMQ(t *testing.T) {
	conn, err := ConnectRabbitMQ()
	require.NoError(t, err)
	require.NotEmpty(t, conn)
}
func TestPublisherRabbitMQ(t *testing.T) {
	conn, err := ConnectRabbitMQ()
	require.NoError(t, err)
	require.NotEmpty(t, conn)
	pub, err := publisher.NewPublisher(conn)
	require.NoError(t, err)
	require.NotEmpty(t, pub)
	message, err := json.Marshal("Hello, world!")
	require.NoError(t, err)
	err = pub.Publish(context.Background(), message, "text/plain")
	require.NoError(t, err)
}
func TestConsumerRabbitMQ(t *testing.T) {
	conn, err := ConnectRabbitMQ()
	require.NoError(t, err)
	require.NotEmpty(t, conn)
	sub, err := consumer.NewConsumer(conn)
	require.NoError(t, err)
	require.NotEmpty(t, sub)
	err = sub.StartConsumer(worker)
	require.NoError(t, err)
}

func ConnectRabbitMQ() (*amqp091.Connection, error) {
	err := utils.LoadFileEnvOnLocal()
	if err != nil {
		return nil, err
	}
	urlRabbitMQ, ok := os.LookupEnv("URL_RABBITMQ")
	if !ok || urlRabbitMQ == "" {
		return nil, errors.New("URL Empty")
	}
	conn, err := NewRabbitMQConn(RabbitMQConnStr(urlRabbitMQ))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
func worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)
		switch delivery.Type {
		case "ordered":
			var mess string
			err := json.Unmarshal(delivery.Body, &mess)
			if err != nil {
				slog.Error("failed to unmarshal message", err)
			}
			slog.Info("MESSAGE::", mess)
			err = delivery.Ack(false)
			if err != nil {
				slog.Error("failed to acknowledge delivery", err)
			}
		default:
			slog.Info("Default")
		}
	}
	slog.Info("Deliveries channel closed")
}

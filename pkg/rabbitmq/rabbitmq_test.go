package rabbitmq

import (
	"errors"
	"os"
	"testing"

	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
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
	publisher, err := publisher.NewPublisher(conn)
	require.NoError(t, err)
	require.NotEmpty(t, publisher)
}

type Test struct {
	sharedkernel.AggregateRoot
	ID      string
	Content string
}
type Test2 struct {
	sharedkernel.DomainEvent
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (t Test2) Identity() string {
	return "Test2"
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

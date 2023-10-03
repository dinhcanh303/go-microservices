package rabbitmq

import (
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

type RabbitMQConnStr string

var ErrCannotConnRabbitMQ = errors.New("cannot connect to RabbitMQ")

func NewRabbitMQConn(rabbitMqURL RabbitMQConnStr) (*amqp.Connection, error) {
	var (
		amqpConn *amqp.Connection
		counts   int64
	)
	for {
		conn, err := amqp.Dial(string(rabbitMqURL))
		if err != nil {
			slog.Error("Failed to connect to RabbitMQ...", err, rabbitMqURL)
			counts++
		} else {
			amqpConn = conn
			break
		}
		if counts > _retryTimes {
			slog.Error("Failed to retry", err)
			return nil, ErrCannotConnRabbitMQ
		}
		slog.Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)
		continue
	}
	slog.Info("📫 connected to RabbitMQ 🎉")
	return amqpConn, nil

}

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/notification/config"
	"github.com/dinhcanh303/go-microservices/pkg/mongodb"
	consumer "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/comsumer"
	"github.com/rabbitmq/amqp091-go"
)

type App struct {
	Cfg      *config.Config
	Mongo    mongodb.MongoDBEngine
	AmqpConn *amqp091.Connection
	Consumer consumer.EventConsumer
}

func New(
	cfg *config.Config,
	mg mongodb.MongoDBEngine,
	amqpConn *amqp091.Connection,
	consumer consumer.EventConsumer,
) *App {
	return &App{
		Cfg:      cfg,
		Mongo:    mg,
		AmqpConn: amqpConn,
		Consumer: consumer,
	}
}

//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/search/config"
	"github.com/dinhcanh303/go-microservices/internal/search/app/router"
	"github.com/dinhcanh303/go-microservices/internal/search/eventhandlers"
	"github.com/dinhcanh303/go-microservices/pkg/elastic"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq"
	consumer "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/comsumer"
	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	rabbitMQConnStr rabbitmq.RabbitMQConnStr,
	elasticSearchConn elastic.ElasticSearchConn,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		rabbitMQFunc,
		elasticSearchFunc,
		router.SearchGRPCServerSet,
		consumer.EventConsumerSet,
		eventhandlers.EventHandlersSet,
	))
}
func rabbitMQFunc(url rabbitmq.RabbitMQConnStr) (*amqp091.Connection, func(), error) {
	conn, err := rabbitmq.NewRabbitMQConn(url)
	if err != nil {
		return nil, nil, err
	}
	return conn, func() { conn.Close() }, nil
}
func elasticSearchFunc(conn elastic.ElasticSearchConn) (elastic.ElasticSearch, error) {
	es, err := elastic.NewElasticSearch(conn)
	if err != nil {
		return nil, err
	}
	return es, err
}

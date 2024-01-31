//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/search/config"
	"github.com/dinhcanh303/go-microservices/internal/search/app/router"
	"github.com/dinhcanh303/go-microservices/internal/search/eventhandlers"
	infrasGRPC "github.com/dinhcanh303/go-microservices/internal/search/infras/grpc"
	"github.com/dinhcanh303/go-microservices/internal/search/usecases/searches"
	"github.com/dinhcanh303/go-microservices/pkg/meili"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq"
	consumer "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/comsumer"
	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	rabbitMQConnStr rabbitmq.RabbitMQConnStr,
	meiliSearchConn meili.MeiliSearchConn,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		rabbitMQFunc,
		meiliSearchFunc,
		infrasGRPC.AuthGRPCClientSet,
		infrasGRPC.GroupGRPCClientSet,
		router.SearchGRPCServerSet,
		consumer.EventConsumerSet,
		searches.UseCaseSet,
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

//	func elasticSearchFunc(conn elastic.ElasticSearchConn) (elastic.ElasticSearch, error) {
//		es, err := elastic.NewElasticSearch(conn)
//		if err != nil {
//			return nil, err
//		}
//		return es, err
//	}
func meiliSearchFunc(conn meili.MeiliSearchConn) meili.MeiliSearch {
	return meili.NewMeiliSearch(conn)
}

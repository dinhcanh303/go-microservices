//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/notification/config"
	"github.com/dinhcanh303/go-microservices/internal/notification/app/router"
	"github.com/dinhcanh303/go-microservices/internal/notification/eventhandlers"
	"github.com/dinhcanh303/go-microservices/internal/notification/infras/repo"
	"github.com/dinhcanh303/go-microservices/internal/notification/usecases/notifications"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	pkgPostgres "github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq"
	consumer "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/consumer"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitApp(
	cfg *config.Config,
	cfg2 *configs.Redis,
	dbConnStr pkgPostgres.DBConnString,
	dbReadConnStr pkgPostgres.DBConnReadString,
	rabbitMQConnStr rabbitmq.RabbitMQConnStr,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		redisEngineFunc,
		rabbitMQFunc,
		router.NotiGRPCServerSet,
		notifications.UseCaseSet,
		eventhandlers.EventHandlersSet,
		repo.RepositoryNotiSet,
		consumer.EventConsumerSet,
	))
}
func dbEngineFunc(url pkgPostgres.DBConnString, urlRead pkgPostgres.DBConnReadString) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(string(url)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
func redisEngineFunc(config *configs.Redis) (redis.RedisEngine, func(), error) {
	redis, err := redis.NewRedisClient(config)
	if err != nil {
		return nil, nil, err
	}
	return redis, func() { redis.Close() }, nil
}
func rabbitMQFunc(url rabbitmq.RabbitMQConnStr) (*amqp091.Connection, func(), error) {
	conn, err := rabbitmq.NewRabbitMQConn(url)
	if err != nil {
		return nil, nil, err
	}
	return conn, func() { conn.Close() }, nil
}

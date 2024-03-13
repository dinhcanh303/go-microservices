//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/like/config"
	"github.com/dinhcanh303/go-microservices/internal/like/app/router"
	"github.com/dinhcanh303/go-microservices/internal/like/infras"
	infrasGRPC "github.com/dinhcanh303/go-microservices/internal/like/infras/grpc"
	"github.com/dinhcanh303/go-microservices/internal/like/infras/repo"
	likesUC "github.com/dinhcanh303/go-microservices/internal/like/usecases/likes"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	cfg2 *configs.Redis,
	dbConnStr postgres.DBConnString,
	dbReadConnStr postgres.DBConnReadString,
	rabbitMQConnStr rabbitmq.RabbitMQConnStr,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		redisEngineFunc,
		rabbitMQFunc,
		publisher.EventPublisherSet,
		router.LikeGRPCServerSet,
		repo.RepositorySet,
		likesUC.UseCaseSet,
		infras.NotiEventPublisherSet,
		infrasGRPC.CommentGRPCClientSet,
		infrasGRPC.PostGRPCClientSet,
	))
}
func dbEngineFunc(url postgres.DBConnString, urlRead postgres.DBConnReadString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url, urlRead)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
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

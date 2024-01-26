//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/app/router"
	infrasGRPC "github.com/dinhcanh303/go-microservices/internal/post/infras/grpc"
	"github.com/dinhcanh303/go-microservices/internal/post/infras/repo"
	postsUC "github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	cfg2 *configs.Redis,
	dbConnStr postgres.DBConnString,
	dbReadConnStr postgres.DBConnReadString,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		redisEngineFunc,
		router.PostGRPCServerSet,
		repo.RepositoryPostSet,
		postsUC.UseCaseSet,
		infrasGRPC.CommentGRPCClientSet,
		infrasGRPC.LikeGRPCClientSet,
		infrasGRPC.UploadGRPCClientSet,
		infrasGRPC.GroupGRPCClientSet,
		infrasGRPC.AuthGRPCClientSet,
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

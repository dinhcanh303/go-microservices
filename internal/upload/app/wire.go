//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/app/handlers"
	"github.com/dinhcanh303/go-microservices/internal/upload/app/router"
	"github.com/dinhcanh303/go-microservices/internal/upload/infras/repo"
	uploadsUC "github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/minio"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	cfg2 *configs.Redis,
	cfgMinio *configs.Minio,
	dbConnStr postgres.DBConnString,
	dbReadConnStr postgres.DBConnReadString,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		minioFunc,
		redisEngineFunc,
		repo.RepositoryUploadSet,
		uploadsUC.UseCaseSet,
		uploadsUC.UseCaseGRPCSet,
		handlers.UploadHandlerSet,
		router.UploadServiceServer,
	))
}
func dbEngineFunc(url postgres.DBConnString, urlRead postgres.DBConnReadString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url, urlRead)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
func minioFunc(cfg *configs.Minio) (minio.MinioService, func(), error) {
	return minio.NewMinio(cfg), func() {}, nil
}
func redisEngineFunc(config *configs.Redis) (redis.RedisEngine, func(), error) {
	redis, err := redis.NewRedisClient(config)
	if err != nil {
		return nil, nil, err
	}
	return redis, func() { redis.Close() }, nil
}

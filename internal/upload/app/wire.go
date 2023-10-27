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
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	cfgMinio *configs.Minio,
	dbConnStr postgres.DBConnString,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		minioFunc,
		repo.RepositoryUploadSet,
		uploadsUC.UseCaseSet,
		uploadsUC.UseCaseGRPCSet,
		handlers.UploadHandlerSet,
		router.UploadServiceServer,
	))
}
func dbEngineFunc(url postgres.DBConnString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
func minioFunc(cfg *configs.Minio) (minio.MinioService, func(), error) {
	return minio.NewMinio(cfg), func() {}, nil
}

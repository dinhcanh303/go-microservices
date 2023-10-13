//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/like/config"
	"github.com/dinhcanh303/go-microservices/internal/like/app/router"
	"github.com/dinhcanh303/go-microservices/internal/like/infras/repo"
	likesUC "github.com/dinhcanh303/go-microservices/internal/like/usecases/likes"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		router.LikeGRPCServerSet,
		repo.RepositorySet,
		likesUC.UseCaseSet,
	))
}
func dbEngineFunc(url postgres.DBConnString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}

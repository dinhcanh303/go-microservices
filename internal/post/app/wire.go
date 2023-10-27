//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/app/router"
	infrasGRPC "github.com/dinhcanh303/go-microservices/internal/post/infras/grpc"
	"github.com/dinhcanh303/go-microservices/internal/post/infras/repo"
	postsUC "github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
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
		router.PostGRPCServerSet,
		repo.RepositoryPostSet,
		postsUC.UseCaseSet,
		infrasGRPC.CommentGRPCClientSet,
		infrasGRPC.LikeGRPCClientSet,
		infrasGRPC.UploadGRPCClientSet,
	))
}
func dbEngineFunc(url postgres.DBConnString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}

//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/app/router"
	infrasGRPC "github.com/dinhcanh303/go-microservices/internal/comment/infras/grpc"
	"github.com/dinhcanh303/go-microservices/internal/comment/infras/repo"
	commentsUC "github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	dbReadConnStr postgres.DBConnReadString,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		router.CommentGRPCServerSet,
		repo.RepositorySet,
		commentsUC.UseCaseSet,
		infrasGRPC.LikeGRPCClientSet,
		infrasGRPC.UploadGRPCClientSet,
	))
}
func dbEngineFunc(url postgres.DBConnString, urlRead postgres.DBConnReadString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url, urlRead)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}

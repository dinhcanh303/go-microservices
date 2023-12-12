//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/app/router"
	"github.com/dinhcanh303/go-microservices/internal/group/infras/repo"
	groupmembersUC "github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	groupsUC "github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
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
		router.GroupGRPCServerSet,
		repo.RepositoryGroupMemberSet,
		groupmembersUC.UseCaseSet,
		repo.RepositoryGroupSet,
		groupsUC.UseCaseSet,
	))
}
func dbEngineFunc(url postgres.DBConnString, urlRead postgres.DBConnReadString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url, urlRead)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}

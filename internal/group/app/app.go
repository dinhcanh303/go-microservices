package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg             *config.Config
	PG              postgres.DBEngine
	UC              groups.UseCase
	UCGroupMember   groupmembers.UseCase
	GroupGRPCServer gen.GroupServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc groups.UseCase,
	ucGroupMember groupmembers.UseCase,
	groupGRPCServer gen.GroupServiceServer) *App {
	return &App{
		Cfg:             cfg,
		UC:              uc,
		UCGroupMember:   ucGroupMember,
		PG:              pg,
		GroupGRPCServer: groupGRPCServer,
	}
}

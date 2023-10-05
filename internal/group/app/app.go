package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg             *config.Config
	PG              postgres.DBEngine
	UC              groups.UseCase
	GroupGRPCServer gen.GroupServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc groups.UseCase,
	groupGRPCServer gen.GroupServiceServer) *App {
	return &App{
		Cfg:             cfg,
		UC:              uc,
		PG:              pg,
		GroupGRPCServer: groupGRPCServer,
	}
}

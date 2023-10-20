package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/attachments"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg             *config.Config
	PG              postgres.DBEngine
	UC              attachments.UseCase
	GroupGRPCServer gen.GroupServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc attachments.UseCase,
	groupGRPCServer gen.GroupServiceServer) *App {
	return &App{
		Cfg:             cfg,
		UC:              uc,
		PG:              pg,
		GroupGRPCServer: groupGRPCServer,
	}
}

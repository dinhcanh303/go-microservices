package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/apikeys"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg             *config.Config
	PG              postgres.DBEngine
	UCApiKey        apikeys.UseCase
	GroupGRPCServer gen.GroupServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	ucApiKey apikeys.UseCase,
	groupGRPCServer gen.GroupServiceServer) *App {
	return &App{
		Cfg:             cfg,
		UCApiKey:        ucApiKey,
		PG:              pg,
		GroupGRPCServer: groupGRPCServer,
	}
}

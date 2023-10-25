package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
)

type App struct {
	Cfg      *config.Config
	CfgMinio *configs.Minio
	PG       postgres.DBEngine
	UC       uploads.UseCase
}

func New(
	cfg *config.Config,
	cfgMinio *configs.Minio,
	pg postgres.DBEngine,
	uc uploads.UseCase,
) *App {
	return &App{
		Cfg:      cfg,
		CfgMinio: cfgMinio,
		UC:       uc,
		PG:       pg,
	}
}

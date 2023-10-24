package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg              *config.Config
	PG               postgres.DBEngine
	UC               uploads.UseCase
	UploadGRPCServer gen.UploadServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc uploads.UseCase,
	uploadGRPCServer gen.UploadServiceServer) *App {
	return &App{
		Cfg:              cfg,
		UC:               uc,
		PG:               pg,
		UploadGRPCServer: uploadGRPCServer,
	}
}

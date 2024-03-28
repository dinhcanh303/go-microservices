package app

import (
	v1 "github.com/dinhcanh303/go-microservices/api/upload/v1"
	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/app/handlers"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
)

type App struct {
	Cfg              *config.Config
	CfgMinio         *configs.Minio
	PG               postgres.DBEngine
	UC               uploads.UseCase
	UcGRPC           uploads.UseCaseGRPC
	Handler          *handlers.UploadHandler
	UploadGRPCServer v1.UploadServiceServer
}

func New(
	cfg *config.Config,
	cfgMinio *configs.Minio,
	pg postgres.DBEngine,
	uc uploads.UseCase,
	UcGRPC uploads.UseCaseGRPC,
	handler *handlers.UploadHandler,
	uploadGRPCServer v1.UploadServiceServer,
) *App {
	return &App{
		Cfg:              cfg,
		CfgMinio:         cfgMinio,
		UC:               uc,
		UcGRPC:           UcGRPC,
		PG:               pg,
		Handler:          handler,
		UploadGRPCServer: uploadGRPCServer,
	}
}

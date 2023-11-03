package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/app/handlers"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg              *config.Config
	CfgMinio         *configs.Minio
	PG               postgres.DBEngine
	UC               uploads.UseCase
	UcGRPC           uploads.UseCaseGRPC
	Handler          *handlers.UploadHandler
	UploadGRPCServer gen.UploadServiceServer
}

func New(
	cfg *config.Config,
	cfgMinio *configs.Minio,
	pg postgres.DBEngine,
	uc uploads.UseCase,
	UcGRPC uploads.UseCaseGRPC,
	handler *handlers.UploadHandler,
	uploadGRPCServer gen.UploadServiceServer,
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

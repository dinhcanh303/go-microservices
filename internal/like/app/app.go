package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/like/config"
	"github.com/dinhcanh303/go-microservices/internal/like/usecases/likes"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg            *config.Config
	PG             postgres.DBEngine
	UC             likes.UseCase
	LikeGRPCServer gen.LikeServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc likes.UseCase,
	likeGRPCServer gen.LikeServiceServer) *App {
	return &App{
		Cfg:            cfg,
		UC:             uc,
		PG:             pg,
		LikeGRPCServer: likeGRPCServer,
	}
}

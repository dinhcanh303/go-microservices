package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg            *config.Config
	PG             postgres.DBEngine
	UC             posts.UseCase
	PostGRPCServer gen.PostServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc posts.UseCase,
	postGRPCServer gen.PostServiceServer) *App {
	return &App{
		Cfg:            cfg,
		UC:             uc,
		PG:             pg,
		PostGRPCServer: postGRPCServer,
	}
}

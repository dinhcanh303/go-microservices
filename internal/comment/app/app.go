package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg               *config.Config
	PG                postgres.DBEngine
	UC                comments.UseCase
	CommentGRPCServer gen.CommentServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc comments.UseCase,
	commentGRPCServer gen.CommentServiceServer) *App {
	return &App{
		Cfg:               cfg,
		UC:                uc,
		PG:                pg,
		CommentGRPCServer: commentGRPCServer,
	}
}

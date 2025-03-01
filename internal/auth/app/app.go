package app

import (
	"context"

	v1 "github.com/dinhcanh303/go-microservices/api/auth/v1"
	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/app/handlers"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/pkg/oauth2"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
)

type App struct {
	Cfg             *config.Config
	PG              postgres.DBEngine
	UC              auth.UseCase
	Listen          auth.ListenTrigger
	AuthGRPCServer  v1.AuthServiceServer
	Handler         *handlers.AuthHandler
	UploadDomainSvc domain.UploadDomainService
	ChangeDBUserPub publisher.EventPublisher
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc auth.UseCase,
	listen auth.ListenTrigger,
	authGRPCServer v1.AuthServiceServer,
	handler *handlers.AuthHandler,
	uploadDomainSvc domain.UploadDomainService,
	changeDBUserPub publisher.EventPublisher) *App {
	return &App{
		Cfg:             cfg,
		UC:              uc,
		PG:              pg,
		AuthGRPCServer:  authGRPCServer,
		Handler:         handler,
		UploadDomainSvc: uploadDomainSvc,
		ChangeDBUserPub: changeDBUserPub,
		Listen:          listen,
	}
}

func (a *App) ListenTrigger(ctx context.Context) {
	a.Listen.ChangeDBUser(ctx)
}

func (a *App) InitOauth2Func() error {
	err := oauth2.InitOAuth()
	if err != nil {
		return err
	}
	return nil
}

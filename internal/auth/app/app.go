package app

import (
	"context"

	v1 "github.com/dinhcanh303/go-microservices/api/auth/v1"
	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
)

type App struct {
	Cfg             *config.Config
	CfgLdap         *configs.Ldap
	PG              postgres.DBEngine
	UC              auth.UseCase
	Listen          auth.ListenTrigger
	AuthGRPCServer  v1.AuthServiceServer
	UploadDomainSvc domain.UploadDomainService
	ChangeDBUserPub publisher.EventPublisher
}

func New(
	cfg *config.Config,
	cfgLdap *configs.Ldap,
	pg postgres.DBEngine,
	uc auth.UseCase,
	listen auth.ListenTrigger,
	authGRPCServer v1.AuthServiceServer,
	uploadDomainSvc domain.UploadDomainService,
	changeDBUserPub publisher.EventPublisher) *App {
	return &App{
		Cfg:             cfg,
		CfgLdap:         cfgLdap,
		UC:              uc,
		PG:              pg,
		AuthGRPCServer:  authGRPCServer,
		UploadDomainSvc: uploadDomainSvc,
		ChangeDBUserPub: changeDBUserPub,
		Listen:          listen,
	}
}

func (a *App) ListenTrigger(ctx context.Context) {
	a.Listen.ChangeDBUser(ctx)
}

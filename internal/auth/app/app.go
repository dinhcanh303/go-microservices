package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg             *config.Config
	CfgLdap         *configs.Ldap
	PG              postgres.DBEngine
	UC              auth.UseCase
	AuthGRPCServer  gen.AuthServiceServer
	UploadDomainSvc domain.UploadDomainService
}

func New(
	cfg *config.Config,
	cfgLdap *configs.Ldap,
	pg postgres.DBEngine,
	uc auth.UseCase,
	authGRPCServer gen.AuthServiceServer,
	uploadDomainSvc domain.UploadDomainService) *App {
	return &App{
		Cfg:             cfg,
		CfgLdap:         cfgLdap,
		UC:              uc,
		PG:              pg,
		AuthGRPCServer:  authGRPCServer,
		UploadDomainSvc: uploadDomainSvc,
	}
}

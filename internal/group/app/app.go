package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
)

type App struct {
	Cfg *config.Config
	PG  postgres.DBEngine
	UC  groups.UseCase
}

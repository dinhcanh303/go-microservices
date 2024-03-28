package app

import (
	"context"

	v1 "github.com/dinhcanh303/go-microservices/api/group/v1"
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
)

type App struct {
	Cfg              *config.Config
	PG               postgres.DBEngine
	UC               groups.UseCase
	UCGroupMember    groupmembers.UseCase
	GroupGRPCServer  v1.GroupServiceServer
	Listen           groups.ListenTrigger
	ChangeDBGroupPub publisher.EventPublisher
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc groups.UseCase,
	ucGroupMember groupmembers.UseCase,
	groupGRPCServer v1.GroupServiceServer,
	listen groups.ListenTrigger,
	changeDBGroupPub publisher.EventPublisher,
) *App {
	return &App{
		Cfg:              cfg,
		UC:               uc,
		UCGroupMember:    ucGroupMember,
		PG:               pg,
		GroupGRPCServer:  groupGRPCServer,
		ChangeDBGroupPub: changeDBGroupPub,
		Listen:           listen,
	}
}
func (a *App) ListenTrigger(ctx context.Context) {
	a.Listen.ChangeDBUser(ctx)
}

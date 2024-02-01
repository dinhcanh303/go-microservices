package app

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg              *config.Config
	PG               postgres.DBEngine
	UC               groups.UseCase
	UCGroupMember    groupmembers.UseCase
	GroupGRPCServer  gen.GroupServiceServer
	Listen           groups.ListenTrigger
	ChangeDBGroupPub publisher.EventPublisher
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc groups.UseCase,
	ucGroupMember groupmembers.UseCase,
	groupGRPCServer gen.GroupServiceServer,
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

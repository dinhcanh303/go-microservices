package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/like/config"
	"github.com/dinhcanh303/go-microservices/internal/like/usecases/likes"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/rabbitmq/amqp091-go"
)

type App struct {
	Cfg            *config.Config
	PG             postgres.DBEngine
	UC             likes.UseCase
	AmqpConn       *amqp091.Connection
	LikeGRPCServer gen.LikeServiceServer
	NotiPub        likes.NotiEventPublisher
	Publisher      publisher.EventPublisher
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc likes.UseCase,
	amqpConn *amqp091.Connection,
	likeGRPCServer gen.LikeServiceServer,
	notiPub likes.NotiEventPublisher,
	publisher publisher.EventPublisher) *App {
	return &App{
		Cfg:            cfg,
		UC:             uc,
		PG:             pg,
		AmqpConn:       amqpConn,
		LikeGRPCServer: likeGRPCServer,
		NotiPub:        notiPub,
		Publisher:      publisher,
	}
}

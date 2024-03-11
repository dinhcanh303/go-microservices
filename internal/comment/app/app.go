package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/rabbitmq/amqp091-go"
)

type App struct {
	Cfg               *config.Config
	PG                postgres.DBEngine
	UC                comments.UseCase
	AmqpConn          *amqp091.Connection
	CommentGRPCServer gen.CommentServiceServer
	LikeDomainSvc     domain.LikeDomainService
	UploadDomainSvc   domain.UploadDomainService
	NotiPub           comments.NotiEventPublisher
	Publisher         publisher.EventPublisher
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc comments.UseCase,
	amqpConn *amqp091.Connection,
	commentGRPCServer gen.CommentServiceServer,
	likeDomainSvc domain.LikeDomainService,
	uploadDomainSvc domain.UploadDomainService,
	notiPub comments.NotiEventPublisher,
	publisher publisher.EventPublisher,
) *App {
	return &App{
		Cfg:               cfg,
		UC:                uc,
		PG:                pg,
		CommentGRPCServer: commentGRPCServer,
		LikeDomainSvc:     likeDomainSvc,
		UploadDomainSvc:   uploadDomainSvc,
		AmqpConn:          amqpConn,
		NotiPub:           notiPub,
		Publisher:         publisher,
	}
}

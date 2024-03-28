package app

import (
	v1 "github.com/dinhcanh303/go-microservices/api/comment/v1"
	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/rabbitmq/amqp091-go"
)

type App struct {
	Cfg               *config.Config
	PG                postgres.DBEngine
	UC                comments.UseCase
	AmqpConn          *amqp091.Connection
	CommentGRPCServer v1.CommentServiceServer
	LikeDomainSvc     domain.LikeDomainService
	UploadDomainSvc   domain.UploadDomainService
	PostDomainSvc     domain.PostDomainService
	NotiPub           comments.NotiEventPublisher
	Publisher         publisher.EventPublisher
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc comments.UseCase,
	amqpConn *amqp091.Connection,
	commentGRPCServer v1.CommentServiceServer,
	likeDomainSvc domain.LikeDomainService,
	uploadDomainSvc domain.UploadDomainService,
	postDomainSvc domain.PostDomainService,
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
		PostDomainSvc:     postDomainSvc,
		AmqpConn:          amqpConn,
		NotiPub:           notiPub,
		Publisher:         publisher,
	}
}

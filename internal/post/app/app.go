package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/rabbitmq/amqp091-go"
)

type App struct {
	Cfg              *config.Config
	PG               postgres.DBEngine
	AMQPConn         *amqp091.Connection
	UC               posts.UseCase
	Publisher        publisher.EventPublisher
	PostGRPCServer   gen.PostServiceServer
	CommentDomainSvc domain.CommentDomainService
	LikeDomainSvc    domain.LikeDomainService
	UploadDomainSvc  domain.UploadDomainService
	NotiPub          posts.NotiEventPublisher
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	amqpConn *amqp091.Connection,
	uc posts.UseCase,
	publisher publisher.EventPublisher,
	postGRPCServer gen.PostServiceServer,
	commentDomainSvc domain.CommentDomainService,
	likeDomainSvc domain.LikeDomainService,
	uploadDomainSvc domain.UploadDomainService,
	notiPub posts.NotiEventPublisher,
) *App {
	return &App{
		Cfg:              cfg,
		UC:               uc,
		PG:               pg,
		AMQPConn:         amqpConn,
		PostGRPCServer:   postGRPCServer,
		CommentDomainSvc: commentDomainSvc,
		LikeDomainSvc:    likeDomainSvc,
		UploadDomainSvc:  uploadDomainSvc,
		Publisher:        publisher,
		NotiPub:          notiPub,
	}
}

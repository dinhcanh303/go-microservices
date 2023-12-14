package infras

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/google/wire"
)

var UserCreatedEventPublisherSet = wire.NewSet(NewUserCreatedEventPublisher)
var UserDeletedEventPublisherSet = wire.NewSet(NewUserDeletedEventPublisher)

type (
	userCreatedEventPublisher struct {
		pub publisher.EventPublisher
	}
	userDeletedEventPublisher struct {
		pub publisher.EventPublisher
	}
)

// Configure implements auth.UserCreatedEventPublisher.
func (p *userCreatedEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

// Publish implements auth.UserCreatedEventPublisher.
func (p *userCreatedEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}

func NewUserCreatedEventPublisher(pub publisher.EventPublisher) auth.UserCreatedEventPublisher {
	return &userCreatedEventPublisher{
		pub: pub,
	}
}

// Configure implements auth.UserDeletedEventPublisher.
func (p *userDeletedEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

// Publish implements auth.UserDeletedEventPublisher.
func (p *userDeletedEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}

func NewUserDeletedEventPublisher(pub publisher.EventPublisher) auth.UserDeletedEventPublisher {
	return &userDeletedEventPublisher{
		pub: pub,
	}
}

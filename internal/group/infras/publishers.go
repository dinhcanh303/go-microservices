package infras

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/google/wire"
)

var GroupCreatedEventPublisherSet = wire.NewSet(NewGroupCreatedEventPublisher)
var GroupDeletedEventPublisherSet = wire.NewSet(NewGroupDeletedEventPublisher)

type (
	groupCreatedEventPublisher struct {
		pub publisher.EventPublisher
	}
	groupDeletedEventPublisher struct {
		pub publisher.EventPublisher
	}
)

func NewGroupCreatedEventPublisher(pub publisher.EventPublisher) groups.GroupCreatedEventPublisher {
	return &groupCreatedEventPublisher{
		pub: pub,
	}
}

// Configure implements groups.GroupEventPublisher.
func (p *groupCreatedEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

// Publish implements groups.GroupEventPublisher.
func (p *groupCreatedEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}

func NewGroupDeletedEventPublisher(pub publisher.EventPublisher) groups.GroupDeletedEventPublisher {
	return &groupDeletedEventPublisher{
		pub: pub,
	}
}

// Configure implements groups.GroupEventPublisher.
func (p *groupDeletedEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

// Publish implements groups.GroupEventPublisher.
func (p *groupDeletedEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}

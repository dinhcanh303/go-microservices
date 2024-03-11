package infras

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/google/wire"
)

type notiEventPublisher struct {
	pub publisher.EventPublisher
}

var _ posts.NotiEventPublisher = (*notiEventPublisher)(nil)

var NotiEventPublisherSet = wire.NewSet(NewNotiEventPublisher)

func NewNotiEventPublisher(pub publisher.EventPublisher) comments.NotiEventPublisher {
	return &notiEventPublisher{
		pub: pub,
	}
}

// Configure implements posts.NotiEventPublisher.
func (p *notiEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

// Publish implements posts.NotiEventPublisher.
func (p *notiEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}

package eventhandlers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	publisher "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
)

type groupUploadEventHandler struct {
	pg       postgres.DBEngine
	groupPub publisher.EventPublisher
}

var _ GroupUploadEventHandler = (*groupUploadEventHandler)(nil)

func NewGroupUploadEventHandler(pg postgres.DBEngine, groupPub publisher.EventPublisher) GroupUploadEventHandler {
	return &groupUploadEventHandler{
		pg:       pg,
		groupPub: groupPub,
	}
}

// Handle implements GroupUploadEventHandler.
func (h *groupUploadEventHandler) Handle(ctx context.Context, e event.GroupUploadBegin) {
	panic("unimplemented")
}

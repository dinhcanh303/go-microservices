package eventhandlers

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	"github.com/dinhcanh303/go-microservices/internal/search/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/elastic"
	"github.com/google/wire"
)

type eventhandlers struct {
	elastic elastic.ElasticSearch
}

var _ EventHandlers = (*eventhandlers)(nil)

var EventHandlersSet = wire.NewSet(NewEventHandlers)

func NewEventHandlers(elastic elastic.ElasticSearch) EventHandlers {
	return &eventhandlers{
		elastic: elastic,
	}

}

// HandleGroupCreated implements EventHandlers.
func (h *eventhandlers) HandleGroupCreated(ctx context.Context, e event.GroupCreated) error {
	groupSearch := domain.NewGroupSearch(e)
	err := h.elastic.Insert(constant.ELASTIC_SEARCH_INDEX, groupSearch, groupSearch.GroupID.String())
	if err != nil {
		slog.Error("insert-into-elastic-search", err)
		return err
	}
	return nil
}

// HandleGroupDeleted implements EventHandlers.
func (h *eventhandlers) HandleGroupDeleted(ctx context.Context, e event.GroupDeleted) error {
	err := h.elastic.Remove(e.GroupID.String(), constant.ELASTIC_SEARCH_INDEX)
	if err != nil {
		slog.Error("remove-from-elastic-search", err)
		return err
	}
	return nil
}

// HandleUserCreated implements EventHandlers.
func (h *eventhandlers) HandleUserCreated(ctx context.Context, e event.UserCreated) error {
	userSearch := domain.NewUserSearch(e)
	err := h.elastic.Insert(constant.ELASTIC_SEARCH_INDEX, userSearch, userSearch.UserID.String())
	if err != nil {
		slog.Error("insert-into-elastic-search", err)
		return err
	}
	return nil
}

// HandleUserDeleted implements EventHandlers.
func (h *eventhandlers) HandleUserDeleted(ctx context.Context, e event.UserDeleted) error {
	err := h.elastic.Remove(e.UserID.String(), constant.ELASTIC_SEARCH_INDEX)
	if err != nil {
		slog.Error("remove-from-elastic-search", err)
		return err
	}
	return nil
}

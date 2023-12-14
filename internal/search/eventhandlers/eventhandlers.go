package eventhandlers

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	"github.com/dinhcanh303/go-microservices/internal/search/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/meili"
	"github.com/google/wire"
)

type eventhandlers struct {
	meili meili.MeiliSearch
}

var _ EventHandlers = (*eventhandlers)(nil)

var EventHandlersSet = wire.NewSet(NewEventHandlers)

func NewEventHandlers(meili meili.MeiliSearch) EventHandlers {
	return &eventhandlers{
		meili: meili,
	}

}

// HandleGroupCreated implements EventHandlers.
func (h *eventhandlers) HandleGroupCreated(ctx context.Context, e event.GroupCreated) error {
	groupSearch := domain.NewGroupSearch(e)
	_, err := h.meili.AddDocuments(constant.MEILI_SEARCH_INDEX, groupSearch)
	if err != nil {
		slog.Error("insert-into-meili-search", err)
		return err
	}
	slog.Info("Insert into meili-search")
	return nil
}

// HandleGroupDeleted implements EventHandlers.
func (h *eventhandlers) HandleGroupDeleted(ctx context.Context, e event.GroupDeleted) error {
	_, err := h.meili.DeleteDocument(constant.MEILI_SEARCH_INDEX, e.ID.String())
	if err != nil {
		slog.Error("remove-from-meili-search", err)
		return err
	}
	return nil
}

// HandleUserCreated implements EventHandlers.
func (h *eventhandlers) HandleUserCreated(ctx context.Context, e event.UserCreated) error {
	userSearch := domain.NewUserSearch(e)
	_, err := h.meili.AddDocuments(constant.MEILI_SEARCH_INDEX, userSearch)
	if err != nil {
		slog.Error("insert-into-meili-search", err)
		return err
	}
	return nil
}

// HandleUserDeleted implements EventHandlers.
func (h *eventhandlers) HandleUserDeleted(ctx context.Context, e event.UserDeleted) error {
	_, err := h.meili.DeleteDocument(constant.MEILI_SEARCH_INDEX, e.ID.String())
	if err != nil {
		slog.Error("remove-from-meili-search", err)
		return err
	}
	return nil
}

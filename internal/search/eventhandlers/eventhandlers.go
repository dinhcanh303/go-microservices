package eventhandlers

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/search/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/meili"
	"github.com/google/wire"
)

type eventhandlers struct {
	meili          meili.MeiliSearch
	authDomainSvc  domain.AuthDomainService
	groupDomainSvc domain.GroupDomainService
}

var _ EventHandlers = (*eventhandlers)(nil)

var EventHandlersSet = wire.NewSet(NewEventHandlers)

func NewEventHandlers(meili meili.MeiliSearch,
	authDomainSvc domain.AuthDomainService,
	groupDomainSvc domain.GroupDomainService,
) EventHandlers {
	return &eventhandlers{
		meili:          meili,
		authDomainSvc:  authDomainSvc,
		groupDomainSvc: groupDomainSvc,
	}

}

// HandleChangeDBGroup implements EventHandlers.
func (e *eventhandlers) HandleChangeDBGroup(ctx context.Context) error {
	_, err := e.meili.DeleteAllDocuments(constant.MeiliSearchDBGroupIndex)
	if err != nil {
		return err
	}
	result, err := e.groupDomainSvc.GetGroups(ctx)
	if err != nil {
		return err
	}
	_, err = e.meili.AddDocuments(constant.MeiliSearchDBUserIndex, result.Groups)
	if err != nil {
		slog.Error("insert-into-meili-search", err)
		return err
	}
	slog.Info("Insert into meili-search")
	return nil
}

// HandleChangeDBUser implements EventHandlers.
func (e *eventhandlers) HandleChangeDBUser(ctx context.Context) error {
	_, err := e.meili.DeleteAllDocuments(constant.MeiliSearchDBUserIndex)
	if err != nil {
		return err
	}
	result, err := e.authDomainSvc.GetUsers(ctx)
	if err != nil {
		return err
	}
	_, err = e.meili.AddDocuments(constant.MeiliSearchDBUserIndex, result.Users)
	if err != nil {
		slog.Error("insert-into-meili-search", err)
		return err
	}
	slog.Info("Insert into meili-search")
	return nil
}

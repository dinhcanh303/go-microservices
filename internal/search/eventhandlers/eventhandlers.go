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
	_, err := e.meili.DeleteAllDocuments(constant.MeiliSearchGroupIndex)
	if err != nil {
		return err
	}
	results, err := e.groupDomainSvc.GetGroups(ctx)
	if err != nil {
		return err
	}
	slog.Info("GROUPS::", results)
	_, err = e.meili.AddDocuments(constant.MeiliSearchGroupIndex, results.GetGroups())
	if err != nil {
		slog.Error("insert-into-meili-search", err)
		return err
	}
	slog.Info("Insert data groups into meili-search")
	return nil
}

// HandleChangeDBUser implements EventHandlers.
func (e *eventhandlers) HandleChangeDBUser(ctx context.Context) error {
	_, err := e.meili.DeleteAllDocuments(constant.MeiliSearchUserIndex)
	if err != nil {
		return err
	}
	results, err := e.authDomainSvc.GetUsers(ctx)
	if err != nil {
		return err
	}
	slog.Info("USERS::", results)
	_, err = e.meili.AddDocuments(constant.MeiliSearchUserIndex, results.GetUsers())
	if err != nil {
		slog.Error("insert-into-meili-search", err)
		return err
	}
	slog.Info("Insert data users into meili-search")
	return nil
}

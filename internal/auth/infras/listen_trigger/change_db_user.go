package listen_trigger

import (
	"context"
	"log/slog"
	"time"

	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	pkgPublisher "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/google/wire"
	"github.com/lib/pq"
)

type changeDBUser struct {
	url             postgres.DBConnString
	changeDBUserPub pkgPublisher.EventPublisher
}

var _ auth.ListenTrigger = (*changeDBUser)(nil)

var ListenTriggerSet = wire.NewSet(NewListenTrigger)

func NewListenTrigger(
	url postgres.DBConnString,
	changeDBUserPub pkgPublisher.EventPublisher) auth.ListenTrigger {
	return &changeDBUser{
		url:             url,
		changeDBUserPub: changeDBUserPub,
	}
}

// ChangeDBUser implements auth.ListenTrigger.
func (tg *changeDBUser) ChangeDBUser(ctx context.Context) {
	reportProblem := func(_ pq.ListenerEventType, err error) {
		if err != nil {
			slog.Error("Report Listener::", err)
		}
	}
	minReconnect := 10 * time.Second
	maxReconnect := time.Minute
	listener := pq.NewListener(string(tg.url), minReconnect, maxReconnect, reportProblem)
	err := listener.Listen(constant.UserChangeEvent)
	if err != nil {
		slog.Error("Listener::", err)
	}
	defer listener.Close()
	slog.Info("entering main loop")
	for {
		waitForNotification(ctx, listener, tg.changeDBUserPub)
	}
}
func waitForNotification(ctx context.Context, l *pq.Listener, changeDBUserPub pkgPublisher.EventPublisher) {
	select {
	case <-l.Notify:
		// slog.Info("received notification, new work available")
		var emptyMessage []byte
		err := changeDBUserPub.Publish(ctx, emptyMessage, "text/plain")
		if err != nil {
			slog.Error("publish message failed", err)
		}
	case <-time.After(90 * time.Second):
		go l.Ping()
		// Check if there's more work available, just in case it takes
		// a while for the Listener to notice connection loss and
		// reconnect.
		// slog.Info("received no work for 90 seconds, checking for new work")
	}
}

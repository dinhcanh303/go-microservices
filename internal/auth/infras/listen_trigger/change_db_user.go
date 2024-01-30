package listen_trigger

import (
	"context"
	"log/slog"
	"time"

	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	pkgPublisher "github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/google/wire"
	"github.com/lib/pq"
	"github.com/pkg/errors"
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
func (tg *changeDBUser) ChangeDBUser() error {
	ctx := context.Background()
	listener := pq.NewListener(string(tg.url),
		10*time.Second,
		time.Minute,
		func(_ pq.ListenerEventType, err error) {
			if err != nil {
				slog.Warn("pq.Listener error:", err)
			}
		})
	err := listener.Listen("user_change_event")
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	defer listener.Close()

	go func() {
		for {
			select {
			case <-listener.Notify:
				// Handle the notification as needed
				var data []byte
				if err := tg.changeDBUserPub.Publish(ctx, data, "text/plain"); err != nil {
					slog.Warn("failed to publish event change database user", err)
				}
			case <-time.After(90 * time.Second):
				go func() {
					err := listener.Ping()
					if err != nil {
						slog.Warn("failed to listener ping", err)
					}
				}()
			}
		}
	}()
	return err
}

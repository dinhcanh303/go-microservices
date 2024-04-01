package nats

import (
	"errors"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

type NatsConnStr string

var ErrCannotConnNats = errors.New("cannot connect to Nats")

func NewNatsConn(natsURL NatsConnStr) (*nats.Conn, error) {
	var (
		natsConn *nats.Conn
		counts   int64
	)
	for {
		nc, err := nats.Connect(string(natsURL))
		if err != nil {
			slog.Error("Failed to connect to Nats...", err, natsURL)
			counts++
		} else {
			natsConn = nc
			break
		}
		if counts > _retryTimes {
			slog.Error("Failed to retry", err)
			return nil, ErrCannotConnNats
		}
		slog.Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)
		continue
	}
	slog.Info("📫 connected to Nats 🎉")
	return natsConn, nil
}

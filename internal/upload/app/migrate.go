//go:build migrate

package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang/glog"
	"golang.org/x/exp/slog"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 5
	_defaultTimeout  = time.Second
)

var (
	_migrationFilePath = "db/migrations"
)

func init() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		glog.Fatalf("migrate: environment variable not declared: PG_URL")
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)
	for attempts > 0 {
		inDocker, ok := os.LookupEnv("IN_DOCKER")
		if !ok || len(inDocker) == 0 {
			glog.Fatalf("migrate: environment variable not declared: IN_DOCKER")
		}
		dir := fmt.Sprintf("file://%s", _migrationFilePath)
		if dockerd, _ := strconv.ParseBool(inDocker); !dockerd {
			cur, _ := os.Getwd()
			dir = fmt.Sprintf("file://%s/%s", filepath.Dir(cur+"/../../.."), _migrationFilePath)
		}
		slog.Info(dir)
		m, err = migrate.New(dir, databaseURL)
		if err == nil {
			break
		}

		slog.Info("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}
	if err != nil {
		glog.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		glog.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("Migrate: no change")
		return
	}
	slog.Info("Migrate: up success")
}

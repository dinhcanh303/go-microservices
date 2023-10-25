package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/app"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/logger"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/slog"
)

func main() {
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("Failed set max process", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed get config", err)
	}
	cfgMinio, err := configs.NewConfigMinio()
	if err != nil {
		slog.Error("Failed get config", err)
	}
	slog.Info("⚡ Init App", "name", cfg.Name, "version", cfg.Version)

	//set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	//integrate Logrus with the slog logger
	logrusHandle := logger.NewLogrusHandler(logrus.StandardLogger())
	slog.New(logrusHandle)
	webPort, ok := os.LookupEnv("WEB_PORT")
	if !ok || webPort == "" {
		glog.Fatalf("web: environment variable not declared: webPort")
	}
	e := echo.New()
	prepareApp(ctx, cancel, cfg, cfgMinio, e)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", webPort)))

}

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, cfgMinio *configs.Minio, echo *echo.Echo) func() {
	_, cleanup, err := app.InitApp(cfg, cfgMinio, postgres.DBConnString(cfg.PG.DsnURL), echo)
	if err != nil {
		slog.Error("Failed init app", err)
		cancel()
		<-ctx.Done()
	}
	return cleanup
}

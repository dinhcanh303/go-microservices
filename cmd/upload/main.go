package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/app"
	"github.com/dinhcanh303/go-microservices/internal/upload/app/handlers"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/logger"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
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
	cfgMinio, err := configs.NewConfigMinio()
	slog.Info("MINIO config::", cfgMinio)
	if err != nil {
		slog.Error("Failed get config", err)
	}
	cfg, err := config.NewConfig()
	slog.Info("MINIO config::", cfg)
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
	e := echo.New()
	prepareApp(ctx, cancel, cfg, cfgMinio, e)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", cfg.HTTP.Port)))

}

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, cfgMinio *configs.Minio, echo *echo.Echo) func() {
	app, cleanup, err := app.InitApp(cfg, cfgMinio, postgres.DBConnString(cfg.PG.DsnURL))
	if err != nil {
		slog.Error("Failed init app", err)
		cancel()
		<-ctx.Done()
	}
	configureRoutes(echo, *app.Handler)
	return cleanup
}
func configureRoutes(echo *echo.Echo, handlers handlers.UploadHandler) {
	echo.GET("/attachments/:id", handlers.GetAttachment)
	echo.POST("/upload", handlers.UploadFile)
	echo.DELETE("/attachments/:id", handlers.DeleteAttachment)
	echo.PUT("/attachments/:id", handlers.UpdateAttachment)
}

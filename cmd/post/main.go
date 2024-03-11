package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/app"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/logger"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
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
	cfgRedis, err := configs.NewConfigRedis()
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

	server := grpc.NewServer()

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()
	cleanup := prepareApp(ctx, cancel, cfg, cfgRedis, server)

	//gRPC Server
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	network := "tcp"
	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("Failed to listen to address", err, "Network", network, "Address", address)
		cancel()
		<-ctx.Done()
	}
	slog.Info("🌏 start server...", "address", address)
	defer func() {
		if err1 := l.Close(); err != nil {
			slog.Error("failed to close", err1, "network", network, "address", address)
			<-ctx.Done()
		}
	}()
	err = server.Serve(l)
	if err != nil {
		slog.Error("failed start gRPC server", err, "network", network, "address", address)
		cancel()
		<-ctx.Done()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case v := <-quit:
		cleanup()
		slog.Info("signal.Notify", v)
	case done := <-ctx.Done():
		cleanup()
		slog.Info("ctx.Done", "app done", done)
	}

}

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, cfgRedis *configs.Redis, server *grpc.Server) func() {
	a, cleanup, err := app.InitApp(cfg, cfgRedis, postgres.DBConnString(cfg.PG.DbURL), postgres.DBConnReadString(cfg.PG.DbRepURL), rabbitmq.RabbitMQConnStr(cfg.RabbitMQ.URL), server)
	if err != nil {
		slog.Error("Failed init app", err)
		cancel()
		<-ctx.Done()
	}
	a.NotiPub.Configure(
		publisher.ExChangeName("post-noti-exchange"),
		publisher.BindingKey("post-noti-routing-key"),
		publisher.MessageTypeName("post-noti"),
	)
	return cleanup
}

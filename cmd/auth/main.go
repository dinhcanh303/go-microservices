package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/app"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/logger"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func main() {
	log := logger.New()
	slog.Info("Build main group started")
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("Failed set max process", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed get config", err)
	}
	cfgLdap, err := configs.NewLdapConfig()
	if err != nil {
		slog.Error("Failed get config Ldap", err)
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
	//
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_validator.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(log.GetZapLogger()),
		)))
	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()
	a, cleanup := prepareApp(ctx, cancel, cfg, cfgRedis, cfgLdap, server)
	//
	//Start Http2 server
	httpMux := http.NewServeMux()
	httpMux.Handle("/api/v1/oauth/google", http.HandlerFunc(a.Handler.GoogleLogin))
	httpMux.Handle("/api/v1/oauth_callback/google", http.HandlerFunc(a.Handler.GoogleCallback))
	http2Server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTP2.Host, cfg.HTTP2.Port),
		Handler: httpMux,
	}
	slog.Info("🌏 start listening http_2...", "address", fmt.Sprintf("%s:%d", cfg.HTTP2.Host, cfg.HTTP2.Port))
	if err := http2Server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to listen and serve", err)
	}
	//Listen change database
	go func() {
		a.ListenTrigger(ctx)
	}()
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

	//
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

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, cfgRedis *configs.Redis, cfgLdap *configs.Ldap, server *grpc.Server) (*app.App, func()) {
	a, cleanup, err := app.InitApp(cfg, cfgRedis, cfgLdap, postgres.DBConnString(cfg.PG.DbURL), postgres.DBConnReadString(cfg.PG.DbRepURL),
		rabbitmq.RabbitMQConnStr(cfg.RabbitMQ.URL), server)
	if err != nil {
		slog.Error("Failed init app", err)
		cancel()
		<-ctx.Done()
	}
	err = a.InitOauth2Func()
	if err != nil {
		slog.Error("Failed init oauth2", err)
		cancel()
		<-ctx.Done()
	}
	a.ChangeDBUserPub.Configure(
		publisher.ExChangeName("search-exchange"),
		publisher.BindingKey("search-routing-key"),
		publisher.MessageTypeName("users-changed"),
	)
	return a, cleanup
}

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dinhcanh303/go-microservices/cmd/gateway/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/app"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/logger"
	"github.com/dinhcanh303/go-microservices/pkg/middleware"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	gatewayRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
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
	slog.Info("⚡ Init App", "name", cfg.Name, "version", cfg.Version)

	//set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	//integrate Logrus with the slog logger
	logrusHandle := logger.NewLogrusHandler(logrus.StandardLogger())
	slog.New(logrusHandle)
	runGrpcServer(ctx, cancel, cfg, cfgLdap)
}

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, cfgLdap *configs.Ldap, server *grpc.Server) (*app.App, func()) {
	app, cleanup, err := app.InitApp(cfg, cfgLdap, postgres.DBConnString(cfg.PG.DsnURL), server)
	if err != nil {
		slog.Error("Failed init app", err)
		cancel()
		<-ctx.Done()
	}
	return app, cleanup
}
func runGrpcServer(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, cfgLdap *configs.Ldap) *app.App {
	server := grpc.NewServer()

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()
	app, cleanup := prepareApp(ctx, cancel, cfg, cfgLdap, server)
	go runGatewayServer(ctx, cancel, cfg, app)
	//gRPC Server
	address := fmt.Sprintf("%s:%d", cfg.AuthHost, cfg.AuthPort)
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
	return app
}
func runGatewayServer(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, app *app.App) {
	mux := http.NewServeMux()
	gw, err := newGateway(ctx, cfg, nil)
	if err != nil {
		slog.Error("failed to create a new gateway", err)
	}
	gwAuth, err := newGatewayAuth(ctx, cfg, nil)
	if err != nil {
		slog.Error("failed to create a new gateway", err)
	}
	mux.Handle("/", gw)
	mux.Handle("/", middleware.AuthMiddleware(gwAuth, app))
	slog.InfoContext(ctx, "Context::")
	//server swagger
	fs := http.FileServer(http.Dir("swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: allowCORS(withLogger(mux)),
	}

	//goroutine
	go func() {
		<-ctx.Done()
		slog.Info("shutting down the http server")

		if err := s.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown http server", err)
		}
	}()

	slog.Info("🌏 start listening...", "address", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))

	if err := s.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to listen and serve", err)
	}
}
func newGatewayAuth(
	ctx context.Context,
	cfg *config.Config,
	opts []gatewayRuntime.ServeMuxOption) (http.Handler, error) {

	groupEndpoint := fmt.Sprintf("%s:%d", cfg.GroupHost, cfg.GroupPort)
	postEndpoint := fmt.Sprintf("%s:%d", cfg.PostHost, cfg.PostPort)
	commentEndpoint := fmt.Sprintf("%s:%d", cfg.CommentHost, cfg.CommentPort)
	likeEndpoint := fmt.Sprintf("%s:%d", cfg.LikeHost, cfg.LikePort)
	uploadEndpoint := fmt.Sprintf("%s:%d", cfg.UploadHost, cfg.UploadPort)

	mux := gatewayRuntime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gen.RegisterGroupServiceHandlerFromEndpoint(ctx, mux, groupEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	err = gen.RegisterPostServiceHandlerFromEndpoint(ctx, mux, postEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	err = gen.RegisterCommentServiceHandlerFromEndpoint(ctx, mux, commentEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	err = gen.RegisterLikeServiceHandlerFromEndpoint(ctx, mux, likeEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	err = gen.RegisterUploadServiceHandlerFromEndpoint(ctx, mux, uploadEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}
func newGateway(
	ctx context.Context,
	cfg *config.Config,
	opts []gatewayRuntime.ServeMuxOption) (http.Handler, error) {

	authEndpoint := fmt.Sprintf("%s:%d", cfg.AuthHost, cfg.AuthPort)

	mux := gatewayRuntime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gen.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, authEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
					preflightHandler(w, r)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	headers := []string{"*"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	slog.Info("preflight request", "http_path", r.URL.Path)
}

func withLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Run Request", "http_method", r.Method, "http_url", r.URL)
		h.ServeHTTP(w, r)
	})
}

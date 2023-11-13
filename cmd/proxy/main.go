package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dinhcanh303/go-microservices/cmd/proxy/config"
	"github.com/dinhcanh303/go-microservices/pkg/logger"
	"github.com/dinhcanh303/go-microservices/pkg/middleware"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/golang/glog"
	gatewayRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGatewayWithMiddleware(
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
func newGatewayWithOutMiddleware(
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

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cfg, err := config.NewConfig()
	if err != nil {
		glog.Fatalf("Config error: %s", err)
	}
	// set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	// integrate Logrus with the slog logger
	slog.New(logger.NewLogrusHandler(logrus.StandardLogger()))
	routerWithMiddleware := routerWithMiddleware(ctx, cfg)
	routerWithoutMiddleware := routerWithoutMiddleware(ctx, cfg)
	mux := http.NewServeMux()
	//with middleware
	mux.Handle("/api", routerWithMiddleware)
	//without middleware
	mux.Handle("/api", routerWithoutMiddleware)
	//server swagger
	// mux.Handle("/", routerWithStatic())

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: mux,
	}
	//goroutine
	go func() {
		<-ctx.Done()
		slog.Info("shutting down the http server")

		if err := s.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown http server", err)
		}
	}()
	slog.Info("start listening...", "address", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err := s.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to listen and serve", err)
	}
}
func routerWithMiddleware(ctx context.Context, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()
	gw, err := newGatewayWithMiddleware(ctx, cfg, nil)
	if err != nil {
		slog.Error("failed to create a new gateway", err)
	}
	mux.Handle("/api/v1", gw)
	return allowCORS(middleware.AuthMiddleware(withLogger(mux)))
}
func routerWithoutMiddleware(ctx context.Context, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()
	gw, err := newGatewayWithOutMiddleware(ctx, cfg, nil)
	if err != nil {
		slog.Error("failed to create a new gateway", err)
	}
	mux.Handle("/api/v1", gw)
	return mux
}
func routerWithStatic() http.Handler {
	mux := http.NewServeMux()
	// Server Swagger
	fs := http.FileServer(http.Dir("swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	return mux
}

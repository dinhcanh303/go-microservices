package metric

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	IncHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func CreateMetrics(cfg configs.Metrics, name string) (Metrics, error) {
	var metric PrometheusMetrics
	metric.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})
	if err := prometheus.Register(metric.HitsTotal); err != nil {
		return nil, err
	}
	metric.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metric.Hits); err != nil {
		return nil, err
	}
	metric.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metric.Times); err != nil {
		return nil, err
	}
	if err := prometheus.Register(prometheus.NewBuildInfoCollector()); err != nil {
		return nil, err
	}
	metricMux := http.NewServeMux()
	metricMux.Handle("/metrics", promhttp.Handler())
	address := fmt.Sprintf("%s:%d", cfg.HostMetric, cfg.PortMetric)
	metricsServer := &http.Server{
		Addr:    address,
		Handler: metricMux,
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-ctx.Done()
		slog.Info("shutting down the http server")
		if err := metricsServer.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown http server", err)
		}
	}()
	slog.Info("🌏 start listening...", "address", address)
	if err := metricsServer.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to listen and serve", err)
	}
	return &metric, nil
}

func (metric *PrometheusMetrics) IncHits(status int, method, path string) {
	metric.HitsTotal.Inc()
	metric.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (metric *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	metric.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}

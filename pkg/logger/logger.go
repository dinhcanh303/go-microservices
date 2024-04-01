package logger

import (
	"context"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

type Logger interface {
	GetZapLogger() *zap.Logger
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
}

type logger struct {
	zapLogger        *zap.Logger
	zapSugaredLogger *zap.SugaredLogger
}

func New() Logger {
	env := "dev"
	var l *zap.Logger
	switch env {
	case "production":
	case "prod":
		l, _ = zap.NewProduction()
	case "development":
	case "dev":
		// l, _ = zap.NewDevelopment()
		l, _ = zap.NewProduction()
	default:
		l = zap.NewNop()
	}
	return &logger{
		zapLogger:        l,
		zapSugaredLogger: l.Sugar(),
	}
}
func (l logger) Log(ctx context.Context, level grpcLogging.Level, msg string, fields ...any) {
	switch level {
	case grpcLogging.LevelDebug:
		l.Debugw(msg, fields)
	case grpcLogging.LevelInfo:
		l.Infow(msg, fields)
	case grpcLogging.LevelWarn:
		l.Warnw(msg, fields)
	case grpcLogging.LevelError:
		l.Errorw(msg, fields)
	}
}
func (l logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

func (l logger) Debug(args ...interface{}) {
	l.zapSugaredLogger.Debug(args...)
}

func (l logger) Debugf(template string, args ...interface{}) {
	l.zapSugaredLogger.Debugf(template, args...)
}

func (l logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Debugw(msg, keysAndValues...)
}
func (l logger) Info(args ...interface{}) {
	l.zapSugaredLogger.Info(args...)
}

func (l logger) Infof(template string, args ...interface{}) {
	l.zapSugaredLogger.Infof(template, args...)
}

func (l logger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Infow(msg, keysAndValues...)
}
func (l logger) Warn(args ...interface{}) {
	l.zapSugaredLogger.Warn(args...)
}

func (l logger) Warnf(template string, args ...interface{}) {
	l.zapSugaredLogger.Warnf(template, args...)
}

func (l logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Warnw(msg, keysAndValues...)
}
func (l logger) Error(args ...interface{}) {
	l.zapSugaredLogger.Error(args...)
}

func (l logger) Errorf(template string, args ...interface{}) {
	l.zapSugaredLogger.Errorf(template, args...)
}

func (l logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Errorw(msg, keysAndValues...)
}
func (l logger) DPanic(args ...interface{}) {
	l.zapSugaredLogger.DPanic(args...)
}

func (l logger) DPanicf(template string, args ...interface{}) {
	l.zapSugaredLogger.DPanicf(template, args...)
}

func (l logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.DPanicw(msg, keysAndValues...)
}

func (l logger) DPanicln(args ...interface{}) {
	l.zapSugaredLogger.DPanic(args...)
}

func (l logger) Panic(args ...interface{}) {
	l.zapSugaredLogger.Panic(args...)
}

func (l logger) Panicf(template string, args ...interface{}) {
	l.zapSugaredLogger.Panicf(template, args...)
}

func (l logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Panicw(msg, keysAndValues...)
}
func (l logger) Fatal(args ...interface{}) {
	l.zapSugaredLogger.Fatal(args...)
}

func (l logger) Fatalf(template string, args ...interface{}) {
	l.zapSugaredLogger.Fatalf(template, args...)
}

func (l logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapSugaredLogger.Fatalw(msg, keysAndValues...)
}

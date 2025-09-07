package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	instance *slog.Logger
)

func Init(app string, level slog.Level) {
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: level},
	)

	instance = slog.New(handler).With(slog.String("app", app))
}

func WithComponent(component string) *slog.Logger {
	return instance.With(slog.String("component", component))
}

func Debug(msg string, args ...any) {
	instance.Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	instance.DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	instance.Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	instance.InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	instance.Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	instance.WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	instance.Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	instance.ErrorContext(ctx, msg, args...)
}

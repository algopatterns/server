package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	defaultLogger *slog.Logger
)

func init() {
	env := os.Getenv("ENVIRONMENT")

	var handler slog.Handler

	if env == "production" {
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}

		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}

		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	defaultLogger = slog.New(handler)
}

func Default() *slog.Logger {
	return defaultLogger
}

func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

// creates a logger with context
func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return defaultLogger
	}

	// extract any logger from context if present
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}

	return defaultLogger
}

// adds logger to context
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// helper type for context key
type loggerKey struct{}

func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func ErrorErr(err error, msg string, args ...any) {
	args = append(args, "error", err)
	defaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
	os.Exit(1)
}

func FatalErr(err error, msg string, args ...any) {
	args = append(args, "error", err)
	defaultLogger.Error(msg, args...)
	os.Exit(1)
}

package logger

import (
	"log/slog"
	"os"
)

const (
	envDevelopment = "development"
)

func NewLogger(env string) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if env == envDevelopment {
		opts.AddSource = true
		opts.Level = slog.LevelDebug
	}
	handler = slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

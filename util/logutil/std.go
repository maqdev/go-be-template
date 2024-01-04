package logutil

import (
	"context"
	"log/slog"
	"os"
)

// NewStdLogger creates a slog logger with specified level
// All errors are sent to stderr & the rest goes to stdout.
func NewStdLogger(level slog.Level) *slog.Logger {
	return slog.New(NewStdHandler(level))
}

// NewStdHandler creates a slog handler with specified level
// All errors are sent to stderr & the rest goes to stdout.
func NewStdHandler(level slog.Level) slog.Handler {
	return stdHandler{
		handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
		errHandler: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		}),
	}
}

type stdHandler struct {
	handler    slog.Handler
	errHandler slog.Handler
}

func (s stdHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return s.errHandler.Enabled(ctx, level)
}

func (s stdHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= slog.LevelError {
		return s.errHandler.Handle(ctx, record)
	}
	return s.handler.Handle(ctx, record)
}

func (s stdHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return stdHandler{
		handler:    s.handler.WithAttrs(attrs),
		errHandler: s.errHandler.WithAttrs(attrs),
	}
}

func (s stdHandler) WithGroup(name string) slog.Handler {
	return stdHandler{
		handler:    s.handler.WithGroup(name),
		errHandler: s.errHandler.WithGroup(name),
	}
}

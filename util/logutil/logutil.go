package logutil

import (
	"context"
	"log/slog"
)

type logUtilKeyType string

const logUtilKey logUtilKeyType = "log"

func Get(ctx context.Context) *slog.Logger {
	l := ctx.Value(logUtilKey)
	if l == nil {
		return slog.Default()
	}
	return l.(*slog.Logger)
}

func With(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, logUtilKey, logger)
}

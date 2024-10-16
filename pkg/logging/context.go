package logging

import (
	"context"
	"log/slog"
)

type ctxKeyLogger struct{}

// FromCtx returns a logger from the context.
func FromCtx(ctx context.Context) *slog.Logger {
	logger, _ := ctx.Value(ctxKeyLogger{}).(*slog.Logger)
	return logger
}

// InjectCtx returns a new context with the logger.
func InjectCtx(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger{}, logger)
}

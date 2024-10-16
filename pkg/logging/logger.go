package logging

import (
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger      *slog.Logger
	defaultLoggerMutex = &sync.Mutex{}
)

func init() {
	defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
}

func Default() *slog.Logger {
	return defaultLogger
}

func SetDefault(logger *slog.Logger) {
	defaultLoggerMutex.Lock()
	defer defaultLoggerMutex.Unlock()

	defaultLogger = logger
}

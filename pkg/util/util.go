package util

import (
	"context"
	"log/slog"
	"os"
)

var DefaultLogger *Logger

// Slog wrapper
type Logger struct {
	*slog.Logger
}

func init() {
	DefaultLogger = NewLogger()
}

func NewLogger() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
	}
}

// Create a new logger that contains values stored in the context
func (l *Logger) FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	traceID, ok := ctx.Value(TraceIDContextKey).(string)
	if !ok {
		return l
	}

	// Log the trace_id if it is present in the context
	newLogger := l.Logger.With(
		slog.String("trace_id", traceID),
	)

	return &Logger{
		Logger: newLogger,
	}
}

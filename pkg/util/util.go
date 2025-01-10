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

// Create a new logger that logs values stored in the context
func (l *Logger) FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	traceID, ok := ctx.Value(TraceIDContextKey).(string)
	if !ok {
		return l
	}

	// If sender_trace_id is not present, only log the trace_id retrieved above
	senderTraceID, ok := ctx.Value(SenderTraceIDContextKey).(string)
	if !ok {
		logger := l.Logger.With(
			slog.String("trace_id", traceID),
		)

		return &Logger{
			Logger: logger,
		}
	}

	// Log sender_trace_id and trace_id
	logger := l.Logger.With(
		slog.String("trace_id", traceID),
		slog.String("sender_trace_id", senderTraceID),
	)

	return &Logger{
		Logger: logger,
	}
}

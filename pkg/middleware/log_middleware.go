package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Wrap the response to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
	// prevent multiple response.WriteHeader clals
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

// Ensure the status code is written, as it is not the case when response is 200 (see https://pkg.go.dev/net/http#ResponseWriter)
func (rw *responseWriter) Write(b []byte) (int, error) {
	// Prevent  WriteHeader to be called multiple times
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// Capture the status code of the request
func (rw *responseWriter) WriteHeader(code int) {
	// Prevent  WriteHeader to be called multiple times
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

/*
This middleware log metadata about the request. It also generates a unique id to identify requests

Logged fields:
  - HTTP Verb
  - Request path
  - Request status code
  - Request trace id
  - Duration
*/

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		ctx := r.Context()
		// retrieve trace id from context
		traceId := ""
		if id, ok := ctx.Value(TraceIDContextKey).(string); ok {
			traceId = id
		}

		// Create a logger with the request metadata
		logger := slog.With(
			"method", r.Method,
			"path", r.URL.Path,
			"trace_id", traceId,
		)

		// Process request...
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Do log when the request is completed
		logger.Info("http request completed",
			"status", wrapped.Status(),
			"duration_ms", duration.Milliseconds(),
		)
	})
}

package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/xaviercrochet/turbo-octo-adventure/pkg/util"
)

// This middleware "enrich" the context with data (i.e. a trace_id)
func RequestContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this field will be used later to identify log entries per request
		traceID := uuid.New().String()

		ctx := r.Context()
		ctx = context.WithValue(ctx, util.TraceIDContextKey, traceID)

		// The value of this header will be present in every log entries
		if senderTraceID := r.Header.Get(util.HeaderSenderTraceID); senderTraceID != "" {
			ctx = context.WithValue(ctx, util.SenderTraceIDContextKey, senderTraceID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

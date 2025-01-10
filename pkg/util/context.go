package util

type contextKey = string

const (
	TraceIDContextKey       contextKey = "trace_id"
	SenderTraceIDContextKey contextKey = "sender_trace_id"
)

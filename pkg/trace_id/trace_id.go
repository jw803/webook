package trace_id

import (
	"context"
)

const TraceIDKey = "traceID"

type TraceId interface {
	GenerateTraceId() string
	GetTraceIDFromContext(ctx context.Context) string
}

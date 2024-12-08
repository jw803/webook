package trace_id

import (
	"context"
	"github.com/google/uuid"
)

type NormalTraceId struct {
}

func NewNormalTraceId() TraceId {
	return &NormalTraceId{}
}

func (t *NormalTraceId) GenerateTraceId() string {
	return uuid.New().String()
}

func (t *NormalTraceId) GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

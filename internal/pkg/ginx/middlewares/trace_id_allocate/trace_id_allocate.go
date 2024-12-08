package trace_id_allocate

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/pkg/trace_id"
)

type TraceIDAllocateMiddlewareBuilder struct {
	traceId trace_id.TraceId
}

func NewTraceIDAllocateHandler(traceId trace_id.TraceId) *TraceIDAllocateMiddlewareBuilder {
	return &TraceIDAllocateMiddlewareBuilder{
		traceId: traceId,
	}
}

func (b *TraceIDAllocateMiddlewareBuilder) IgnorePaths(path string) *TraceIDAllocateMiddlewareBuilder {
	return b
}

func (b *TraceIDAllocateMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := b.traceId.GenerateTraceId()
		ctx.Set(trace_id.TraceIDKey, traceId)
		ctx.Next()
	}
}

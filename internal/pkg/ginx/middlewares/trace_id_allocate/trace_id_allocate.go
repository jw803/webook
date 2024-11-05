package trace_id_allocate

import (
	cst_traceid "bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/trace"
	"github.com/gin-gonic/gin"
)

type TraceIDAllocateMiddlewareBuilder struct{}

func NewTraceIDAllocateHandler() *TraceIDAllocateMiddlewareBuilder {
	return &TraceIDAllocateMiddlewareBuilder{}
}

func (l *TraceIDAllocateMiddlewareBuilder) IgnorePaths(path string) *TraceIDAllocateMiddlewareBuilder {
	return l
}

func (m *TraceIDAllocateMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := cst_traceid.NewTraceID()
		ctx.Set(cst_traceid.TraceIDKey, traceID)
		ctx.Next()
	}
}

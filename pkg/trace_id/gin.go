package trace_id

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GinTraceId struct {
}

func NewGinTraceId() TraceId {
	return &GinTraceId{}
}

func (t *GinTraceId) GenerateTraceId() string {
	return uuid.New().String()
}

func (t *GinTraceId) GetTraceIDFromContext(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		traceID := t.getTraceIDFromGinContext(ginCtx)
		return traceID
	}
	fmt.Println("it is not gin context")
	return ""
}

func (t *GinTraceId) getTraceIDFromGinContext(ctx *gin.Context) string {
	if traceID, ok := ctx.Get(TraceIDKey); ok {
		if id, isString := traceID.(string); isString {
			return id
		}
	}
	return ""
}

package loggerx

import (
	"context"
	"github.com/jw803/webook/pkg/trace_id"
	"go.uber.org/zap"
)

const (
	alertLevelKey = "alert_level"
	traceIdKey    = "trace_id"
)

var _ Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	traceId trace_id.TraceId
	l       *zap.Logger
}

func NewZapLogger(traceId trace_id.TraceId, l *zap.Logger) *ZapLogger {
	return &ZapLogger{
		traceId: traceId,
		l:       l,
	}
}

func (z *ZapLogger) Debug(ctx context.Context, msg string, args ...Field) {
	z.l.Debug(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) Info(ctx context.Context, msg string, args ...Field) {
	z.l.Info(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) Warn(ctx context.Context, msg string, args ...Field) {
	z.l.Warn(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) Error(ctx context.Context, msg string, args ...Field) {
	z.l.Error(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) toZapFields(ctx context.Context, args []Field) []zap.Field {
	traceId := z.traceId.GetTraceIDFromContext(ctx)
	res := make([]zap.Field, 0, len(args))
	res = append(res, zap.Any(traceIdKey, traceId))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}

func (z *ZapLogger) P0(ctx context.Context, msg string, args ...Field) {
	args = append(args, Field{
		Key:   alertLevelKey,
		Value: P0,
	})
	z.l.Error(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) P1(ctx context.Context, msg string, args ...Field) {
	args = append(args, Field{
		Key:   alertLevelKey,
		Value: P1,
	})
	z.l.Error(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) P2(ctx context.Context, msg string, args ...Field) {
	args = append(args, Field{
		Key:   alertLevelKey,
		Value: P2,
	})
	z.l.Warn(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) P3(ctx context.Context, msg string, args ...Field) {
	args = append(args, Field{
		Key:   alertLevelKey,
		Value: P3,
	})
	z.l.Warn(msg, z.toZapFields(ctx, args)...)
}

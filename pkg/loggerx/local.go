package loggerx

import (
	"context"
	"github.com/jw803/webook/pkg/trace_id"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	localAlertLevelKey = "alert_level"
	localTraceIdKey    = "trace_id"
)

var _ Logger = (*LocalLogger)(nil)

type LocalLogger struct {
	traceId trace_id.TraceId
	l       *zap.Logger
}

func NewLocalLogger() *LocalLogger {
	traceId := trace_id.NewNormalTraceId()
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &LocalLogger{
		traceId: traceId,
		l:       l,
	}
}

func (z *LocalLogger) Debug(ctx context.Context, msg string, args ...Field) {
	z.l.Debug(msg, z.toZapFields(ctx, args)...)
}

func (z *LocalLogger) Info(ctx context.Context, msg string, args ...Field) {
	z.l.Info(msg, z.toZapFields(ctx, args)...)
}

func (z *LocalLogger) Warn(ctx context.Context, msg string, args ...Field) {
	z.l.Warn(msg, z.toZapFields(ctx, args)...)
}

func (z *LocalLogger) Error(ctx context.Context, msg string, args ...Field) {
	z.l.Error(msg, z.toZapFields(ctx, args)...)
}

func (z *LocalLogger) toZapFields(ctx context.Context, args []Field) []zap.Field {
	traceId := z.traceId.GetTraceIDFromContext(ctx)
	res := make([]zap.Field, 0, len(args))
	res = append(res, zap.Any(localTraceIdKey, traceId))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}

func (z *LocalLogger) P0(ctx context.Context, msg string, args ...Field) {
	args = append(args, Field{
		Key:   localAlertLevelKey,
		Value: P0,
	})
	z.l.Error(msg, z.toZapFields(ctx, args)...)
}

func (z *LocalLogger) P1(ctx context.Context, msg string, args ...Field) {
	args = append(args, Field{
		Key:   localAlertLevelKey,
		Value: P1,
	})
	z.l.Error(msg, z.toZapFields(ctx, args)...)
}

func (z *LocalLogger) P2(ctx context.Context, msg string, args ...Field) {
	args = append(args, Field{
		Key:   localAlertLevelKey,
		Value: P2,
	})
	z.l.Warn(msg, z.toZapFields(ctx, args)...)
}

func (z *LocalLogger) P3(ctx context.Context, msg string, args ...Field) {
	args = append(args, Field{
		Key:   localAlertLevelKey,
		Value: P3,
	})
	z.l.Warn(msg, z.toZapFields(ctx, args)...)
}

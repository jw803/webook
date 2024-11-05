package logging

import (
	"context"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/errorx"
	cst_traceid "bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/trace"

	"go.uber.org/zap"
)

const (
	alertLevelKey = "alert_level"
	errorTypeKey  = "error_type"
)

var _ Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	l  *zap.Logger
	sl *zap.SugaredLogger
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	l.Sugar()
	return &ZapLogger{
		l:  l,
		sl: l.Sugar(),
	}
}

func (z *ZapLogger) Debug(ctx context.Context, msg string, args ...Field) {
	z.l.Debug(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) Info(ctx context.Context, msg string, args ...Field) {
	z.l.Info(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) Warn(ctx context.Context, alertLevel AlertLevel, errorType errorx.ErrorType, msg string, args ...Field) {
	args = z.processDefaultFields(args, alertLevel, errorType)
	z.l.Warn(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) Error(ctx context.Context, alertLevel AlertLevel, errorType errorx.ErrorType, msg string, args ...Field) {
	args = z.processDefaultFields(args, alertLevel, errorType)
	z.l.Error(msg, z.toZapFields(ctx, args)...)
}

func (z *ZapLogger) toZapFields(ctx context.Context, args []Field) []zap.Field {
	traceId := cst_traceid.GetTraceIDFromContext(ctx)
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	if traceId != "" {
		res = append(res, zap.Any(string(cst_traceid.TraceIDKey), traceId))
	}
	return res
}

func (z *ZapLogger) processDefaultFields(args []Field, alertLevel AlertLevel, errorType errorx.ErrorType) []Field {
	args = append(args, Field{
		Key:   alertLevelKey,
		Value: alertLevel,
	})
	args = append(args, Field{
		Key:   errorTypeKey,
		Value: errorType,
	})
	return args
}

func (z *ZapLogger) P0(ctx context.Context, errorType errorx.ErrorType, msg string, args ...Field) {
	z.Error(ctx, P0, errorType, msg, args...)
}

func (z *ZapLogger) P1(ctx context.Context, errorType errorx.ErrorType, msg string, args ...Field) {
	z.Error(ctx, P1, errorType, msg, args...)
}

func (z *ZapLogger) P2(ctx context.Context, errorType errorx.ErrorType, msg string, args ...Field) {
	z.Warn(ctx, P2, errorType, msg, args...)
}

func (z *ZapLogger) P3(ctx context.Context, errorType errorx.ErrorType, msg string, args ...Field) {
	z.Warn(ctx, P3, errorType, msg, args...)
}

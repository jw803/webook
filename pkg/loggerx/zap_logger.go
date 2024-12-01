package loggerx

import (
	"context"
	"go.uber.org/zap"
)

const (
	alertLevelKey = "alert_level"
	errorTypeKey  = "error_type"
)

var _ Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	l  *zap.Logger
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	l.Sugar()
	return &ZapLogger{
		l:  l,
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
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}

func (z *ZapLogger) processDefaultFields(args []Field, alertLevel AlertLevel) []Field {
	args = append(args, Field{
		Key:   alertLevelKey,
		Value: alertLevel,
	})
	return args
}

func (z *ZapLogger) P0(ctx context.Context, msg string, args ...Field) {
	z.Error(ctx, msg, z.processDefaultFields(args, P0)...)
}

func (z *ZapLogger) P1(ctx context.Context, msg string, args ...Field) {
	z.Error(ctx, msg, z.processDefaultFields(args, P1)...)
}

func (z *ZapLogger) P2(ctx context.Context, msg string, args ...Field) {
	z.Warn(ctx, msg, z.processDefaultFields(args, P2)...)
}

func (z *ZapLogger) P3(ctx context.Context, msg string, args ...Field) {
	z.Warn(ctx, msg z.processDefaultFields(args, P3)...)
}

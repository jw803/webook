package loggerx

import (
	"context"
)

type NoOpLogger struct {
}

func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}

func (n NoOpLogger) Debug(ctx context.Context, msg string, args ...Field) {
}

func (n NoOpLogger) Info(ctx context.Context, msg string, args ...Field) {
}

func (n NoOpLogger) Warn(ctx context.Context, msg string, args ...Field) {
}

func (n NoOpLogger) Error(ctx context.Context, msg string, args ...Field) {
}

func (n NoOpLogger) P0(ctx context.Context, msg string, args ...Field) {
}

func (n NoOpLogger) P1(ctx context.Context, msg string, args ...Field) {
}

func (n NoOpLogger) P2(ctx context.Context, msg string, args ...Field) {
}

func (n NoOpLogger) P3(ctx context.Context, msg string, args ...Field) {
}

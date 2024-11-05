package logging

import (
	"context"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/errorx"
)

type NoOpLogger struct {
}

func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}

func (n *NoOpLogger) Debug(ctx context.Context, msg string, args ...Field) {
}

func (n *NoOpLogger) Info(ctx context.Context, msg string, args ...Field) {
}

func (n *NoOpLogger) Warn(ctx context.Context, alertLevel AlertLevel, errorType errorx.ErrorType, msg string, args ...Field) {
}

func (n *NoOpLogger) Error(ctx context.Context, alertLevel AlertLevel, errorType errorx.ErrorType, msg string, args ...Field) {
}

func (n *NoOpLogger) P0(ctx context.Context, errorType errorx.ErrorType, msg string, args ...Field) {
}

func (n *NoOpLogger) P1(ctx context.Context, errorType errorx.ErrorType, msg string, args ...Field) {
}

func (n *NoOpLogger) P2(ctx context.Context, errorType errorx.ErrorType, msg string, args ...Field) {
}

func (n *NoOpLogger) P3(ctx context.Context, errorType errorx.ErrorType, msg string, args ...Field) {
}

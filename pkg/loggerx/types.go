package logging

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...Field)
	Info(ctx context.Context, msg string, args ...Field)
	Warn(ctx context.Context, msg string, args ...Field)
	Error(ctx context.Context, msg string, args ...Field)

	P0(ctx context.Context, msg string, args ...Field)
	P1(ctx context.Context, msg string, args ...Field)
	P2(ctx context.Context, msg string, args ...Field)
	P3(ctx context.Context, msg string, args ...Field)
}

type Field struct {
	Key   string
	Value any
}

type AlertLevel string

const (
	P0 AlertLevel = "P0"
	P1 AlertLevel = "P1"
	P2 AlertLevel = "P2"
	P3 AlertLevel = "P3"
)

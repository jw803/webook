package test_ioc

import (
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/jw803/webook/pkg/trace_id"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLog() loggerx.Logger {
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	l := zap.New(
		zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	logger := loggerx.NewZapLogger(trace_id.NewNormalTraceId(), l)
	return logger
}

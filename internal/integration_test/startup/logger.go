package startup

import (
	"os"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() logging.Logger {
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	l := zap.New(
		zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	logger := logging.NewZapLogger(l)
	return logger
}

package startup

import (
	"github.com/jw803/webook/config"
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/jw803/webook/pkg/trace_id"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitLogger() loggerx.Logger {
	env := config.Get().AppEnv

	var logger loggerx.Logger
	switch env {
	case "test":
		logger = loggerx.NewNoOpLogger()
	default:
		cfg := zap.NewDevelopmentConfig()
		err := viper.UnmarshalKey("log", &cfg)
		if err != nil {
			panic(err)
		}
		l, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		logger = loggerx.NewZapLogger(trace_id.NewNormalTraceId(), l)
	}
	return logger
}

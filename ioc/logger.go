package ioc

import (
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/jw803/webook/pkg/trace_id"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitLogger() loggerx.Logger {
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return loggerx.NewZapLogger(trace_id.NewNormalTraceId(), l)
}

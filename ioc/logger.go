package ioc

import (
	"github.com/jw803/webook/pkg/loggerx"
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
	return loggerx.NewZapLogger(l)
}

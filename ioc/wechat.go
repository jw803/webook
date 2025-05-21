package ioc

import (
	"github.com/jw803/webook/internal/service/oauth2/wechat"
	"github.com/jw803/webook/pkg/loggerx"
	"os"
)

func InitWechatService(l loggerx.Logger) wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("env variable WECHAT_APP_ID is not found")
	}
	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("env variable WECHAT_APP_SECRET is not found")
	}
	return wechat.NewService(appID, appSecret, l)
}

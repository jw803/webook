package startup

import (
	"github.com/jw803/webook/internal/service/oauth2/wechat"
	"github.com/jw803/webook/pkg/loggerx"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService(l loggerx.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}

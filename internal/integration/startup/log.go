package startup

import (
	"github.com/jw803/webook/pkg/loggerx"
)

func InitLog() loggerx.LoggerV1 {
	return loggerx.NewNoOpLogger()
}

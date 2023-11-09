package startup

import (
	"github.com/jw803/webook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}

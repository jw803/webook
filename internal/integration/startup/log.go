package startup

import (
	"github.com/jw803/webook/pkg/loggerx"
)

func InitLog() loggerx.Logger {
	return loggerx.NewNoOpLogger()
}

package startup

import (
	"github.com/jw803/webook/pkg/timex"
	"time"
)

var NowTimeString = "2024-01-01T00:00:00+00:00"

func NewNowFunc(nowTimeString string, timeZone ...string) timex.NowFunc {
	if nowTimeString == "" {
		nowTimeString = "2024-01-01T00:00:00Z"
	}
	var TZ = "Asia/Taipei"
	if len(timeZone) > 0 {
		TZ = timeZone[0]
	}
	return func() time.Time {
		Location, err := time.LoadLocation(TZ)
		if err != nil {
			Location, _ = time.LoadLocation(TZ)
		}
		fakeTime, _ := time.ParseInLocation(time.RFC3339, nowTimeString, Location)
		return fakeTime
	}
}

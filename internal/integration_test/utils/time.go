package utils

import (
	"time"
)

var InitialTimeString = "2024-01-01T08:00:00+08:00"
var InitialTime = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

type NowFunc = func() time.Time

func NewNowFunc(nowTimeString string, timeZone string) NowFunc {
	if nowTimeString == "" {
		nowTimeString = "2024-01-01T00:00:00Z"
	}
	if timeZone == "" {
		timeZone = "Asia/Taipei"
	}
	return func() time.Time {
		Location, err := time.LoadLocation(timeZone)
		if err != nil {
			Location, _ = time.LoadLocation(timeZone)
		}
		fakeTime, _ := time.ParseInLocation(time.RFC3339, nowTimeString, Location)
		return fakeTime
	}
}

package startup

import "time"

var NowTimeString = "2024-01-01T00:00:00+00:00"

func NewNowFunc(timeZone string, timeString ...string) func() time.Time {
	var nowTimeString string
	if len(timeString) == 0 {
		nowTimeString = NowTimeString
	} else {
		nowTimeString = timeString[0]
	}
	Location, err := time.LoadLocation(timeZone)
	if err != nil {
		Location = nil
	}
	return func() time.Time {
		fakeTime, _ := time.ParseInLocation(time.RFC3339, nowTimeString, Location)
		return fakeTime
	}
}

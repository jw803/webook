package timex

import "time"

type NowFunc func() time.Time

func NewNowFunc() NowFunc {
	return func() time.Time {
		return time.Now()
	}
}

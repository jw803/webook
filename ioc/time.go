package ioc

import "github.com/jw803/webook/pkg/timex"

func NewNowFunc() timex.NowFunc {
	return timex.NewNowFunc()
}

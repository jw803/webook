package def

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jw803/webook/config"
)

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

// 只有在 WECHAT_STATE_KEY 沒設定時才會用到，讓本地開發不必額外配置也能跑起來；
// 正式環境務必透過環境變數覆蓋。
const defaultStateKey = "k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"

func StateKey() []byte {
	if k := config.Get().WechatStateKey; k != "" {
		return []byte(k)
	}
	return []byte(defaultStateKey)
}

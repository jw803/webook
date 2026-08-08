package jwtx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jw803/webook/config"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// 只有在對應的環境變數（APP_SECRET_KEY / JWT_REFRESH_TOKEN_KEY）沒設定時才會用到，
// 讓本地開發不必額外配置也能跑起來；正式環境務必透過環境變數覆蓋。
const (
	defaultAtKey = "95osj3fUD7fo0mlYdDbncXz4VD2igvf0"
	defaultRtKey = "95osj3fUD7fo0mlYdDbncXz4VD2igvfx"
)

func AtKey() []byte {
	if k := config.Get().APPSecretKey; k != "" {
		return []byte(k)
	}
	return []byte(defaultAtKey)
}

func RtKey() []byte {
	if k := config.Get().JWTRefreshTokenKey; k != "" {
		return []byte(k)
	}
	return []byte(defaultRtKey)
}

type RedisJWTHandler struct {
	cmd redis.Cmdable
	l   loggerx.Logger
}

func NewRedisHandler(cmd redis.Cmdable, l loggerx.Logger) JWTHandler {
	return &RedisJWTHandler{
		cmd: cmd,
		l:   l,
	}
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(ctx, uid, ssid)
	return err
}

func (h *RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey())
	if err != nil {
		h.l.Error(ctx, "failed to sign jwt string", loggerx.Error(err))
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context, claim *UserClaims) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claim.Ssid), "", time.Hour*24*7).Err()
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	val, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		// Exists 會依據是否存在去進行+1計數，不存在的話結果會是0
		if val == 0 {
			return nil
		}
		return errorx.WithCode(errcode.ErrSessionInvalid, "invalid session")
	default:
		return err
	}
}

func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	// 我现在用 JWT 来校验
	tokenHeader := ctx.GetHeader("Authorization")
	// segs := strings.SplitN(tokenHeader, " ", 2)
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey())
	if err != nil {
		h.l.Error(ctx, "failed to sign jwt string", loggerx.Error(err))
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jw803/webook/internal/pkg/errcode"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
)

// RefreshToken 可以同时刷新长短 token，用 redis 来记录是否有效，即 refresh_token 是一次性的
// 参考登录校验部分，比较 User-Agent 来增强安全性
func (h *UserHandler) RefreshToken(ctx *gin.Context) (any, error) {
	// 只有这个接口，拿出来的才是 refresh_token，其它地方都是 access token
	refreshToken := h.ExtractToken(ctx)
	var rc jwtx.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return jwtx.RtKey, nil
	})
	if err != nil || !token.Valid {
		h.l.P3(ctx, "invalid refresh token")
		return nil, errorx.WithCode(errcode.ErrTokenInvalid, "invalid refresh token")
	}

	err = h.CheckSession(ctx, rc.Ssid)
	if errorx.IsCode(err, errcode.ErrSessionInvalid) {
		// 已经退出登录
		h.l.P3(ctx, "invalid refresh token")
		return nil, errorx.WithCode(errcode.ErrTokenInvalid, "invalid refresh token")
	}
	if err != nil {
		//redis 有问题
		h.l.P1(ctx, "failed to check session", loggerx.Error(err))
		return nil, errorx.WithCode(errcode.ErrSystem, "failed to check session")
	}

	// 搞个新的 access_token
	err = h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		h.l.P1(ctx, "failed to set jwt token", loggerx.Error(err))
		return nil, errorx.WithCode(errcode.ErrSystem, "failed to set jwt token")
	}

	return nil, nil
}

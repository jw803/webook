package wechat

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jw803/webook/internal/interface/web/wechat/def"
	"github.com/jw803/webook/pkg/loggerx"
	uuid "github.com/lithammer/shortuuid/v4"
)

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) (any, error) {
	state := uuid.New()
	val := o.svc.AuthURL(ctx, state)
	err := o.setStateCookie(ctx, state)
	if err != nil {
		o.l.P1(ctx, "failed to set state cookie", loggerx.Error(err))
		return nil, err
	}
	return val, nil
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := def.StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)
	if err != nil {
		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr,
		600, "/oauth2/wechat/callback",
		"", false, true)
	return nil
}

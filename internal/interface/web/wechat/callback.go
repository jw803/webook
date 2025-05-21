package wechat

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
)

type CallBackQuery struct {
	Code string `form:"code"`
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context, query CallBackQuery) (any, error) {
	// 你校验不校验都可以
	code := query.Code
	// state := ctx.Query("state")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if errorx.IsCode(err, errcode.ErrWeChatVerificationCodeInvalid) {
		o.l.P3(ctx, "wechat verification code invalid", loggerx.Error(err))
		return nil, err
	}
	if err != nil {
		o.l.P1(ctx, "failed to verify wechat verification code", loggerx.Error(err))
		return nil, err
	}
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		o.l.P1(ctx, "failed to find or create user by wechat", loggerx.Error(err))
		return nil, err
	}
	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		o.l.P1(ctx, "failed to set login token", loggerx.Error(err))
		return nil, err
	}

	return nil, nil
}

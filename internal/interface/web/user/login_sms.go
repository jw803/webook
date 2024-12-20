package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
)

type loginSMSReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (h *UserHandler) LoginSMS(ctx *gin.Context, req loginSMSReq) (any, error) {
	ok, err := h.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	if !ok {
		h.l.P3(ctx, "invalid sms code")
		return nil, errorx.WithCode(errcode.ErrSMSCodeInvalid, "invalid sms code")
	}

	// 我这个手机号，会不会是一个新用户呢？
	// 这样子
	user, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		h.l.P1(ctx, "failed to find and create user", loggerx.Error(err))
		return nil, errorx.WithCode(errcode.ErrSystem, "failed to find and create user")
	}

	if err = h.SetLoginToken(ctx, user.Id); err != nil {
		// 记录日志
		h.l.P1(ctx, "failed to set login token", loggerx.Error(err))
		return nil, errorx.WithCode(errcode.ErrSystem, "failed to set login token")
	}

	return nil, nil
}

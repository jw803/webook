package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
)

type loginSMSSendCodeReq struct {
	Phone string `json:"phone"`
}

func (h *UserHandler) SendLoginSMSCode(ctx *gin.Context, req loginSMSSendCodeReq) (any, error) {
	err := h.codeSvc.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		return nil, nil
	case errorx.IsCode(err, errcode.ErrSMSCodeSendTooFrequently):
		h.l.P3(ctx, "client send sms code too frequently", loggerx.Error(err))
		return nil, err
	default:
		h.l.P1(ctx, "failed to send sms code", loggerx.Error(err))
		return nil, err
	}
}

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
)

type userEditReq struct {
	Id       int64  `json:"id"`
	NickName string `json:"nickName" validate:"required,max=20"`
	Birthday string `json:"birthday" validate:"required,datetime=2006-01-02"`
	Intro    string `json:"intro" validate:"required,max=60"`
}

func (u *UserHandler) Edit(ctx *gin.Context, req userEditReq) (any, error) {
	err := u.svc.EditExtraInfo(ctx, domain.User{
		Id:       req.Id,
		NickName: req.NickName,
		Birthday: req.Birthday,
		Intro:    req.Intro,
	})
	if errorx.IsCode(err, errcode.ErrUserNotFound) {
		u.l.P2(ctx, "It may be an update by a malicious user")
		return nil, err
	}
	if err != nil {
		u.l.P1(ctx, "It may be an update by a malicious user")
		return nil, err
	}
	return nil, nil
}

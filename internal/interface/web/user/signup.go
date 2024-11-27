package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
)

type signUpReq struct {
	Email           string `json:"email" binding:"required,email"`
	ConfirmPassword string `json:"confirmPassword"`
	Password        string `json:"password"`
}

func (u *UserHandler) SignUp(ctx *gin.Context, req signUpReq) (any, error) {
	if req.ConfirmPassword != req.Password {
		u.l.P3(ctx, "the password and the confirmation password do not match")
		return nil, errorx.WithCode(errcode.ErrPasswordNotMatch,
			"the password and the confirmation password do not match")
	}
	ok, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		u.l.P1(ctx, "failed to validate password", loggerx.Error(err))
		return nil, errorx.WithCode(errcode.ErrSystem, "failed to validate password")
	}
	if !ok {
		u.l.P3(ctx, "invalid password")
		return nil, errorx.WithCode(errcode.ErrInvalidPassword, "invalid password")
	}

	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errorx.IsCode(err, errcode.ErrUserDuplicated) {
		u.l.P3(ctx, "duplicate email", loggerx.Error(err))
		return nil, err
	}
	if err != nil {
		u.l.P1(ctx, "failed to sign up", loggerx.Error(err))
		return nil, err
	}

	return nil, nil
}

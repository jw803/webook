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

func (h *UserHandler) SignUp(ctx *gin.Context, req signUpReq) (any, error) {
	if req.ConfirmPassword != req.Password {
		h.l.P3(ctx, "the password and the confirmation password do not match")
		return nil, errorx.WithCode(errcode.ErrPasswordNotMatch,
			"the password and the confirmation password do not match")
	}
	ok, err := h.passwordExp.MatchString(req.Password)
	if err != nil {
		h.l.P1(ctx, "failed to validate password", loggerx.Error(err))
		return nil, errorx.WithCode(errcode.ErrSystem, "failed to validate password")
	}
	if !ok {
		h.l.P3(ctx, "invalid password")
		return nil, errorx.WithCode(errcode.ErrInvalidPassword, "invalid password")
	}

	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errorx.IsCode(err, errcode.ErrDuplicateEmailSignUp) {
		h.l.P3(ctx, "duplicate email", loggerx.Error(err))
		return nil, err
	}
	if err != nil {
		h.l.P1(ctx, "failed to sign up", loggerx.Error(err))
		return nil, errorx.WithCode(errcode.ErrSystem, err.Error())
	}

	return nil, nil
}

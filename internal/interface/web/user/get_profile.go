package user

import (
	"github.com/gin-gonic/gin"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
)

type Profile struct {
	Email    string
	Phone    string
	Nickname string
	Birthday string
}

func (h *UserHandler) ProfileJWT(ctx *gin.Context, uc *jwtx.UserClaims) (any, error) {
	u, err := h.svc.Profile(ctx, uc.Uid)
	if err != nil {
		h.l.P1(ctx, "failed to get user profile")
	}
	return Profile{
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.NickName,
		Birthday: u.Birthday,
	}, nil
}

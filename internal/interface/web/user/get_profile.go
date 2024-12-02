package user

import (
	"github.com/gin-gonic/gin"
	ijwt "github.com/jw803/webook/internal/interface/web/jwtx"
)

type Profile struct {
	Email    string
	Phone    string
	Nickname string
	Birthday string
}

func (h *UserHandler) ProfileJWT(ctx *gin.Context, uc *ijwt.UserClaims) (any, error) {
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

package user

import (
	"github.com/gin-gonic/gin"
	ijwt "github.com/jw803/webook/internal/interface/web/jwtx"
)

func (u *UserHandler) ProfileJWT(ctx *gin.Context, claims *ijwt.UserClaims) (any, error) {
	println(claims.Uid)
	return "你的 profile", nil
}

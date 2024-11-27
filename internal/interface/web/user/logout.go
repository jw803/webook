package user

import (
	"github.com/gin-gonic/gin"
	ijwt "github.com/jw803/webook/internal/interface/web/jwtx"
)

func (u *UserHandler) LogoutJWT(ctx *gin.Context, claim *ijwt.UserClaims) (any, error) {
	err := u.ClearToken(ctx, claim)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

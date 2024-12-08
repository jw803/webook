package user

import (
	"github.com/gin-gonic/gin"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
)

func (h *UserHandler) LogoutJWT(ctx *gin.Context, claim *jwtx.UserClaims) (any, error) {
	err := h.ClearToken(ctx, claim)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

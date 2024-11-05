package jwt_handler

import (
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/ginx"
	"github.com/gin-gonic/gin"
)

type JWTHandler interface {
	ExtractTokenString(ctx *gin.Context) string
}

type ShoplineTokenClaims = ginx.ShoplineClaims

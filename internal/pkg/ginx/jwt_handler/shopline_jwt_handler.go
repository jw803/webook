package jwt_handler

import (
	"github.com/gin-gonic/gin"
)

type shoplineJWTHandler struct{}

func NewShoplineJWTHandler() JWTHandler {
	return &shoplineJWTHandler{}
}

func (h *shoplineJWTHandler) ExtractTokenString(ctx *gin.Context) string {
	authToken := ctx.Request.Header.Get("x-shopline-session-token")
	return authToken
}

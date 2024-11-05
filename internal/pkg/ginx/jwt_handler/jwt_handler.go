package jwt_handler

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type jWTHandler struct{}

func NewJWTHandler() JWTHandler {
	return &jWTHandler{}
}

func (h *jWTHandler) ExtractTokenString(ctx *gin.Context) string {
	authValue := ctx.GetHeader("Authorization")
	if authValue == "" {
		return ""
	}
	authSegments := strings.SplitN(authValue, " ", 2)
	if len(authSegments) != 2 {
		return ""
	}
	return authSegments[1]
}

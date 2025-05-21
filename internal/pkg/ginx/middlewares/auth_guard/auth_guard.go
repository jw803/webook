package auth_guard

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/v2/sets"
	"github.com/emirpasic/gods/v2/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jw803/webook/config"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/internal/pkg/ginx"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
)

type Claims struct {
	jwt.RegisteredClaims
}

type JWTAuthzMiddlewareBuilder struct {
	publicPaths sets.Set[string]
	jwtx.JWTHandler
	l loggerx.Logger
}

func NewJWTAuthzHandler(jwtHandler jwtx.JWTHandler, l loggerx.Logger) JWTAuthzMiddlewareBuilder {
	return JWTAuthzMiddlewareBuilder{
		publicPaths: hashset.New[string](),
		JWTHandler:  jwtHandler,
		l:           l,
	}
}

func (b JWTAuthzMiddlewareBuilder) IgnorePaths(path string) JWTAuthzMiddlewareBuilder {
	b.publicPaths.Add(path)
	return b
}

func (b JWTAuthzMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if b.publicPaths.Contains(ctx.Request.URL.Path) {
			return
		}
		authToken := b.ExtractToken(ctx)
		c := Claims{}
		token, err := jwt.ParseWithClaims(authToken, &c, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				b.l.P3(ctx, "invalid token")
				return nil, fmt.Errorf("%w: %s", errors.New(""), token.Header["alg"])
			}
			return []byte(config.Get().APPSecretKey), nil
		})
		if err != nil || !token.Valid {
			ginx.WriteResponse(ctx, errorx.WithCode(errcode.ErrTokenInvalid, "invalid token"), "")
			ctx.Abort()
			return
		}
		ctx.Set("claim", c)
		ctx.Next()
	}
}

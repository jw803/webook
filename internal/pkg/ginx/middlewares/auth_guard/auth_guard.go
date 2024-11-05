package auth_guard

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/emirpasic/gods/v2/sets"
	"github.com/emirpasic/gods/v2/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/config"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/ginx"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/ginx/jwt_handler"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/ginx/response"
)

type JWTAuthzMiddlewareBuilder struct {
	publicPaths sets.Set[string]
	jwt_handler.JWTHandler
}

func NewJWTAuthzHandler(jwtHandler jwt_handler.JWTHandler) *JWTAuthzMiddlewareBuilder {
	return &JWTAuthzMiddlewareBuilder{
		publicPaths: hashset.New[string](),
		JWTHandler:  jwtHandler,
	}
}

func (l *JWTAuthzMiddlewareBuilder) IgnorePaths(path string) *JWTAuthzMiddlewareBuilder {
	l.publicPaths.Add(path)
	return l
}

func (m *JWTAuthzMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if m.publicPaths.Contains(ctx.Request.URL.Path) {
			return
		}

		authToken := m.ExtractTokenString(ctx)
		sc := ginx.ShoplineClaims{}

		token, err := jwt.ParseWithClaims(authToken, &sc, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %s", errors.New(""), token.Header["alg"])
			}
			return []byte(config.Get().ClientSecretKey), nil
		})
		if err != nil || !token.Valid {
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			ctx.Abort()
			return
		}
		if time.Unix(int64(sc.Exp), 0).Before(time.Now()) && authToken != config.Get().SmokeTestToken {
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			ctx.Abort()
			return
		}

		ctx.Set("shopline", sc)
		ctx.Next()
	}
}

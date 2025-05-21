package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jw803/webook/internal/interface/web/wechat/def"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/internal/pkg/ginx"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
)

type VerifyStateMiddlewareBuilder struct {
	CookieStateKey string
	QueryStateKey  string
	key            []byte
	l              loggerx.Logger
}

func NewVerifyStateMiddlewareBuilder(l loggerx.Logger) *VerifyStateMiddlewareBuilder {
	return &VerifyStateMiddlewareBuilder{
		key: []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"),
		l:   l,
	}
}

func (b *VerifyStateMiddlewareBuilder) SetCookieStateKey(cookieStateKey string) *VerifyStateMiddlewareBuilder {
	b.CookieStateKey = cookieStateKey
	return b
}

func (b *VerifyStateMiddlewareBuilder) SetQueryStateKey(queryStateKey string) *VerifyStateMiddlewareBuilder {
	b.QueryStateKey = queryStateKey
	return b
}

func (b *VerifyStateMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state := ctx.Query(b.QueryStateKey)
		ck, err := ctx.Cookie(b.CookieStateKey)
		if err != nil {
			b.l.P3(ctx, "failed to get cookie from context", loggerx.Error(err))
			ginx.WriteResponse(ctx, errorx.WithCode(errcode.ErrCookieMissing, err.Error()), nil)
			ctx.Abort()
		}
		var sc def.StateClaims
		_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
			return b.key, nil
		})
		if err != nil {
			b.l.P3(ctx, "failed to parse token", loggerx.Error(err))
			ginx.WriteResponse(ctx, errorx.WithCode(errcode.ErrTokenInvalid, err.Error()), nil)
			ctx.Abort()
		}
		if state != sc.State {
			// state 不匹配，有人搞你
			b.l.P3(ctx, "state mismatch", loggerx.String("query_state", state),
				loggerx.String("cookie_state", sc.State))
			ginx.WriteResponse(ctx, errorx.WithCode(errcode.ErrCookieMissing, "state mismatch"), nil)
			ctx.Abort()
		}
		ctx.Next()
	}
}

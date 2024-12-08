package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	articlehdl "github.com/jw803/webook/internal/interface/web/article"
	userhdl "github.com/jw803/webook/internal/interface/web/user"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/pkg/ginx/middlewares/access_log"
	"github.com/jw803/webook/internal/pkg/ginx/middlewares/auth_guard"
	"github.com/jw803/webook/pkg/loggerx"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *userhdl.UserHandler, articleHdl *articlehdl.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	return server
}

func GinMiddlewares(jwtHdl jwtx.Handler, l loggerx.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		auth_guard.NewJWTAuthzHandler(jwtHdl, l).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/login").
			Build(),
		access_log.NewMiddlewareBuilder(func(ctx context.Context, al *access_log.AccessLog) {
			// 设置为 DEBUG 级别
			l.Debug(ctx, "GIN 收到请求", loggerx.Field{
				Key:   "req",
				Value: al,
			})
		}).AllowReqBody(true).AllowRespBody().Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}

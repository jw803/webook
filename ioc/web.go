package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/web"
	ijwt "github.com/jw803/webook/internal/web/jwt"
	"github.com/jw803/webook/internal/web/middleware"
	"github.com/jw803/webook/pkg/ginx/middlewares/accesslog"
	"github.com/jw803/webook/pkg/logger"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	oauth2WechatHdl *web.OAuth2WechatHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	return server
}

func GinMiddlewares(redisClient redis.Cmdable,
	jwtHdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/login").
			Build(),
		accesslog.NewMiddlewareBuilder(func(ctx context.Context, al *accesslog.AccessLog) {
			// 设置为 DEBUG 级别
			l.Debug("GIN 收到请求", logger.Field{
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

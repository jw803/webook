//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/internal/web"
	ijwt "github.com/jw803/webook/internal/web/jwt"
	"github.com/jw803/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		dao.NewGORMUserDAO,

		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,

		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		service.NewUserService,
		service.NewSMSCodeService,

		// 直接基于内存实现
		ioc.InitSmsMemoryService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ioc.NewWechatHandlerConfig,
		ijwt.NewRedisHandler,

		ioc.InitWebServer,
		ioc.GinMiddlewares,
	)
	// 這邊隨便
	return new(gin.Engine)
}

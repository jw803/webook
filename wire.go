//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jw803/webook/internal/repository"
	repository2 "github.com/jw803/webook/internal/repository/article"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	dao2 "github.com/jw803/webook/internal/repository/dao/article"
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
		dao2.NewGORMArticleDAO,

		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,

		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,
		repository2.NewArticleRepository,

		service.NewUserService,
		service.NewSMSCodeService,
		service.NewArticleService,

		// 直接基于内存实现
		ioc.InitSmsMemoryService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ioc.NewWechatHandlerConfig,
		ijwt.NewRedisHandler,

		ioc.InitWebServer,
		ioc.GinMiddlewares,
	)
	// 這邊隨便
	return new(gin.Engine)
}

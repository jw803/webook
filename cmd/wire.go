//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jw803/webook/internal/interface/web/article"
	"github.com/jw803/webook/internal/interface/web/user"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/repository"
	repository2 "github.com/jw803/webook/internal/repository/article"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	dao2 "github.com/jw803/webook/internal/repository/dao/article"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/ioc"
)

var eventProvider = wire.NewSet(
	ioc.NewConsumers,
)

var webProvider = wire.NewSet(
	jwtx.NewRedisHandler,
	user.NewUserHandler,
	article.NewArticleHandler,

	ioc.GinMiddlewares,
	ioc.InitWebServer,
)

func InitApp() *App {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		dao.NewGORMUserDAO,
		dao2.NewGORMArticleDAO,

		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisArticleCache,

		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,
		repository2.NewArticleRepository,

		service.NewUserService,
		service.NewSMSCodeService,
		service.NewArticleService,
		ioc.InitSmsMemoryService,

		eventProvider,
		webProvider,

		wire.Struct(new(App), "*"),
	)
	// 這邊隨便
	return new(App)
}

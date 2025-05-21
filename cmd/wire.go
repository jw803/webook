//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	article2 "github.com/jw803/webook/internal/interface/event/article"
	"github.com/jw803/webook/internal/interface/web/article"
	"github.com/jw803/webook/internal/interface/web/user"
	"github.com/jw803/webook/internal/interface/web/wechat"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/repository"
	repository2 "github.com/jw803/webook/internal/repository/article"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	dao2 "github.com/jw803/webook/internal/repository/dao/article"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/ioc"
	"github.com/jw803/webook/pkg/samarax"
)

var thirdProvider = wire.NewSet(
	ioc.InitDB, ioc.InitRedis, ioc.InitSaramaClient,
	ioc.InitLogger,
	ioc.NewNowFunc,
	ioc.NewUuidFn,
)

var eventProvider = wire.NewSet(
	ioc.NewConsumers,
	article2.NewInteractiveReadEventConsumer,
	samarax.NewSamaraxBaseHandler,
)

var webProvider = wire.NewSet(
	ioc.InitWeberver,
	ioc.GinMiddlewares,
	jwtx.NewRedisHandler,
	user.NewUserHandler,
	article.NewArticleHandler,
	wechat.NewOAuth2WechatHandler,
)

func InitAPP() *App {
	wire.Build(
		thirdProvider,

		dao.NewGORMInteractiveDAO,
		dao.NewGORMUserDAO,
		dao2.NewGORMArticleDAO,

		cache.NewInteractiveRedisCache,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisArticleCache,

		repository.NewCachedInteractiveRepository,
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,
		repository2.NewArticleRepository,

		service.NewUserService,
		service.NewSMSCodeService,
		service.NewArticleService,
		ioc.InitWechatService,
		ioc.InitSmsMemoryService,

		eventProvider,
		webProvider,

		wire.Struct(new(App), "*"),
	)
	// 這邊隨便
	return new(App)
}

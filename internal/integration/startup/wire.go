//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/internal/repository/article"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	article2 "github.com/jw803/webook/internal/repository/dao/article"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/internal/web"
	ijwt "github.com/jw803/webook/internal/web/jwt"
	"github.com/jw803/webook/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewGORMUserDAO,
	cache.NewRedisUserCache,
	repository.NewCachedUserRepository,
	service.NewUserService)

//go:generate wire
func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		article2.NewGORMArticleDao,

		cache.NewRedisCodeCache,
		repository.NewCachedCodeRepository,
		article.NewArticleRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSmsMemoryService,

		// 指定啥也不干的 wechat service
		InitPhantomWechatService,
		service.NewSMSCodeService,
		service.NewArticleService,
		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ioc.NewWechatHandlerConfig,
		ijwt.NewRedisHandler,

		// gin 的中间件
		ioc.GinMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider,
		article2.NewGORMArticleDao,
		article.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler)
	return new(web.ArticleHandler)
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil)
}

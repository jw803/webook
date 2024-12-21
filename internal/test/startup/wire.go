//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	article2 "github.com/jw803/webook/internal/interface/web/article"
	"github.com/jw803/webook/internal/interface/web/user"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/internal/repository/article"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	articleDao "github.com/jw803/webook/internal/repository/dao/article"
	"github.com/jw803/webook/internal/service"
	ioc2 "github.com/jw803/webook/internal/test/test_ioc"
	"github.com/jw803/webook/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, ioc2.InitLog)
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
		articleDao.NewGORMArticleDAO,

		cache.NewRedisArticleCache,
		cache.NewRedisCodeCache,
		repository.NewCachedCodeRepository,
		article.NewArticleRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSmsMemoryService,

		// 指定啥也不干的 wechat service
		service.NewSMSCodeService,
		service.NewArticleService,

		// handler 部分
		jwtx.NewRedisHandler,
		user.NewUserHandler,
		article2.NewArticleHandler,

		// gin 的中间件
		ioc.GinMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler(dao articleDao.ArticleDAO) *article2.ArticleHandler {
	wire.Build(thirdProvider,
		cache.NewRedisArticleCache,
		article.NewArticleRepository,
		service.NewArticleService,
		article2.NewArticleHandler)
	return new(article2.ArticleHandler)
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil)
}

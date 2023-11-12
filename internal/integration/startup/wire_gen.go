// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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
	"github.com/jw803/webook/internal/web/jwt"
	"github.com/jw803/webook/ioc"
)

// Injectors from wire.go:

//go:generate wire
func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	handler := jwt.NewRedisHandler(cmdable)
	loggerV1 := InitLog()
	v := ioc.GinMiddlewares(cmdable, handler, loggerV1)
	gormDB := InitTestDB()
	userDAO := dao.NewGORMUserDAO(gormDB)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	smsService := ioc.InitSmsMemoryService(cmdable)
	codeService := service.NewSMSCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := InitPhantomWechatService(loggerV1)
	wechatHandlerConfig := ioc.NewWechatHandlerConfig()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler, wechatHandlerConfig)
	articleDao := article2.NewGORMArticleDao(gormDB)
	articleRepository := article.NewArticleRepository(articleDao)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	return engine
}

func InitArticleHandler() *web.ArticleHandler {
	gormDB := InitTestDB()
	articleDao := article2.NewGORMArticleDao(gormDB)
	articleRepository := article.NewArticleRepository(articleDao)
	articleService := service.NewArticleService(articleRepository)
	loggerV1 := InitLog()
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	return articleHandler
}

func InitUserSvc() service.UserService {
	gormDB := InitTestDB()
	userDAO := dao.NewGORMUserDAO(gormDB)
	cmdable := InitRedis()
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	return userService
}

// wire.go:

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)

var userSvcProvider = wire.NewSet(dao.NewGORMUserDAO, cache.NewRedisUserCache, repository.NewCachedUserRepository, service.NewUserService)

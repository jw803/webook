// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	article3 "github.com/jw803/webook/internal/interface/web/article"
	"github.com/jw803/webook/internal/interface/web/user"
	"github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/repository"
	article2 "github.com/jw803/webook/internal/repository/article"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	"github.com/jw803/webook/internal/repository/dao/article"
	"github.com/jw803/webook/internal/service"
	ioc2 "github.com/jw803/webook/internal/test/test_ioc"
	"github.com/jw803/webook/ioc"
)

// Injectors from wire.go:

//go:generate wire
func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	handler := jwtx.NewRedisHandler(cmdable)
	logger := ioc2.InitLog()
	v := ioc.GinMiddlewares(handler, logger)
	gormDB := InitTestDB()
	userDAO := dao.NewGORMUserDAO(gormDB)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	smsService := ioc.InitSmsMemoryService(cmdable)
	codeService := service.NewSMSCodeService(codeRepository, smsService)
	userHandler := user.NewUserHandler(userService, codeService, handler, logger)
	articleDAO := article.NewGORMArticleDAO(gormDB)
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := article2.NewArticleRepository(articleDAO, articleCache)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := article3.NewArticleHandler(articleService, logger)
	engine := ioc.InitWebServer(v, userHandler, articleHandler)
	return engine
}

func InitArticleHandler(dao2 article.ArticleDAO) *article3.ArticleHandler {
	cmdable := InitRedis()
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := article2.NewArticleRepository(dao2, articleCache)
	articleService := service.NewArticleService(articleRepository)
	logger := ioc2.InitLog()
	articleHandler := article3.NewArticleHandler(articleService, logger)
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

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, ioc2.InitLog)

var userSvcProvider = wire.NewSet(dao.NewGORMUserDAO, cache.NewRedisUserCache, repository.NewCachedUserRepository, service.NewUserService)

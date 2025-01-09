//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/jw803/webook/internal/interface/web/article"
	"github.com/jw803/webook/internal/interface/web/user"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/ioc"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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

func InitApp(
	db *gorm.DB,
	redis redis.Cmdable,
) *user.UserHandler {
	wire.Build(
		ioc.InitLogger,

		dao.NewGORMUserDAO,
		cache.NewUserCacheV1,
		cache.NewRedisCodeCache,

		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		service.NewUserService,
		service.NewSMSCodeService,
		ioc.InitSmsMemoryService,

		jwtx.NewRedisHandler,
		user.NewUserHandler,
	)
	// 這邊隨便
	return new(user.UserHandler)
}

//go:build wireinject
// +build wireinject

package startup

import (
	"github.com/google/wire"
	article2 "github.com/jw803/webook/internal/interface/event/article"
	"github.com/jw803/webook/internal/interface/web/user"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/internal/repository/cache"
	"github.com/jw803/webook/internal/repository/dao"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/ioc"
	"github.com/jw803/webook/pkg/samarax"
)

var thirdProvider = wire.NewSet(
	ioc.InitDB, ioc.InitRedis, ioc.InitSaramaClient,
	ioc.InitLogger,
	ioc.NewUuidFn,
)

var eventProvider = wire.NewSet(
	ioc.NewConsumers,
	article2.NewInteractiveReadEventConsumer,
	samarax.NewSamaraxBaseHandler,
)

func InitUserhandler(
	userDAO dao.UserDAO,
) *user.UserHandler {
	wire.Build(
		thirdProvider,

		cache.NewRedisUserCache,
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

package service

import (
	"context"
	"errors"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	EditExtraInfo(ctx context.Context, u domain.User) error
}

type userService struct {
	repo repository.UserRepository
	l    loggerx.Logger
}

// NewUserService 我用的人，只管用，怎么初始化我不管，我一点都不关心如何初始化
func NewUserService(repo repository.UserRepository, l loggerx.Logger) UserService {
	return &userService{
		repo: repo,
		l:    l,
	}
}

//func NewUserServiceV1(f repository.UserRepositoryFactory) UserService {
//	return &userService{
//		// 我在这里，不同的 factory，会创建出来不同实现
//		codeRepo: f.NewRepo(),
//	}
//}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if errorx.IsCode(err, errcode.ErrUserNotFound) {
		return domain.User{}, errorx.WithCode(errcode.ErrInvalidUserNameOrPassword, err.Error())
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码了
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		svc.l.Error(ctx, "the password user inputted is incorrect", loggerx.Error(err))
		return domain.User{}, errorx.WithCode(errcode.ErrInvalidUserNameOrPassword, "the password user inputted is incorrect")
	}
	return u, nil
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	err = svc.repo.Create(ctx, u)
	if errorx.IsCode(err, errcode.ErrUserDuplicated) {
		return errorx.WithCode(errcode.ErrDuplicateEmailSignUp, "duplicate mail signup")
	} else {
		return err
	}
	return nil
}

func (svc *userService) FindOrCreate(ctx context.Context,
	phone string) (domain.User, error) {
	// 这时候，这个地方要怎么办？
	// 这个叫做快路径
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 要判断，有咩有这个用户
	if !errorx.IsCode(err, errcode.ErrUserNotFound) {
		// 绝大部分请求进来这里
		// nil 会进来这里
		// 不为 ErrUserNotFound 的也会进来这里
		return u, err
	}
	// 这里，把 phone 脱敏之后打出来

	//loggerxx.Logger.Info("用户未注册", zap.String("phone", phone))
	// 在系统资源不足，触发降级之后，不执行慢路径了
	//if ctx.Value("降级") == "true" {
	//	return domain.User{}, errors.New("系统降级了")
	//}
	// 这个叫做慢路径
	// 你明确知道，没有这个用户
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && !errorx.IsCode(err, errcode.ErrUserDuplicated) {
		return u, err
	}
	// 因为这里会遇到主从延迟的问题
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context,
	info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenID)
	if !errorx.IsCode(err, errcode.ErrUserNotFound) {
		return u, err
	}
	u = domain.User{
		WechatInfo: info,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && !errorx.IsCode(err, errcode.ErrUserDuplicated) {
		return u, err
	}
	// 因为这里会遇到主从延迟的问题
	return svc.repo.FindByWechat(ctx, info.OpenID)
}

func (svc *userService) Profile(ctx context.Context,
	id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	return u, err
}

func PathsDownGrade(ctx context.Context, quick, slow func()) {
	quick()
	if ctx.Value("降级") == "true" {
		return
	}
	slow()
}

func (svc *userService) EditExtraInfo(ctx context.Context, u domain.User) error {
	if err := svc.repo.EditExtraInfo(ctx, u); err != nil {
		return err
	}
	return nil
}

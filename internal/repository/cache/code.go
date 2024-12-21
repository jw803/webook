package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
	l      loggerx.Logger
}

// NewCodeCacheGoBestPractice Go 的最佳实践是返回具体类型
func NewCodeCacheGoBestPractice(client redis.Cmdable, l loggerx.Logger) *RedisCodeCache {
	return &RedisCodeCache{
		client: client,
		l:      l,
	}
}

func NewRedisCodeCache(client redis.Cmdable, l loggerx.Logger) CodeCache {
	return &RedisCodeCache{
		client: client,
		l:      l,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		c.l.Error(ctx, "failed to set code in redis", loggerx.Error(err))
		return err
	}
	switch res {
	case 0:
		// 毫无问题
		return nil
	case -1:
		// 你要在对应的告警系统里面配置，
		// 比如说规则，一分钟内出现超过100次 WARN，你就告警
		c.l.Error(ctx, "send code too frequently")
		return errorx.WithCode(errcode.ErrSMSCodeSendTooFrequently, "verification code is being sent too frequently")
	//case -2:
	//	return
	default:
		// 系统错误
		c.l.Error(ctx, "unexpected key ttl, it should not be permanent")
		return errorx.WithCode(errcode.ErrRedis, "redis error")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 正常来说，如果频繁出现这个错误，你就要告警，因为有人搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
		//default:
		//	return false, ErrUnknownForCode
	}
	return false, ErrUnknownForCode
}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

// LocalCodeCache 假如说你要切换这个，你是不是得把 lua 脚本的逻辑，在这里再写一遍？
type LocalCodeCache struct {
	client redis.Cmdable
}

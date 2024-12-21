package cache

import (
	"context"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/internal/repository/cache/redismocks"
	"github.com/jw803/webook/internal/test/test_ioc"
	"github.com/jw803/webook/internal/test/test_model"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable

		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr     error
		WantErrCode int
	}{
		{
			name: "successfully set the verification code to Redis",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mockRedis := redismocks.NewMockCmdable(ctrl)
				redisRes := redis.NewCmd(context.Background())
				redisRes.SetVal(int64(0))
				mockRedis.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:152"}, "123").Return(redisRes)
				return mockRedis
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "152",
			code:  "123",

			wantErr: nil,
		},
		{
			name: "redis error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mockRedis := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errorx.WithCode(test_model.ErrRedis, "mock redis error"))
				mockRedis.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:152"}, "123").Return(res)
				return mockRedis
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "152",
			code:  "123",

			wantErr: errorx.WithCode(test_model.ErrRedis, "mock redis error"),
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:152"},
					[]any{"123456"},
				).Return(res)
				return cmd
			},

			ctx:   context.Background(),
			biz:   "login",
			phone: "152",
			code:  "123456",

			wantErr: errorx.WithCode(errcode.ErrSMSCodeSendTooFrequently, "verification code is being sent too frequently"),
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mockRedis := redismocks.NewMockCmdable(ctrl)
				redisRes := redis.NewCmd(context.Background())
				redisRes.SetVal(int64(2))
				mockRedis.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:152"}, "123").Return(redisRes)
				return mockRedis
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "152",
			code:  "123",

			wantErr: errorx.WithCode(errcode.ErrRedis, "redis error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			redisMock := tc.mock(ctrl)
			codeCache := NewRedisCodeCache(redisMock, test_ioc.InitLog())
			err := codeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.True(t, errorx.IsEqual(tc.wantErr, err))
		})
	}
}

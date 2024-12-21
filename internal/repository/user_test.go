package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository/cache"
	cachemocks "github.com/jw803/webook/internal/repository/cache/mocks"
	"github.com/jw803/webook/internal/repository/dao"
	daomocks "github.com/jw803/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_UserRepository_FindById(t *testing.T) {
	nowTSMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowTSMs)

	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		id  int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "success, but cache miss",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := daomocks.NewMockUserDAO(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.Users{
					Id: 123,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "this is password",
					Phone: sql.NullString{
						String: "15212345678",
						Valid:  true,
					},
					Ctime: nowTSMs,
					Utime: nowTSMs,
				}, nil)
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, errors.New("cache error"))
				uc.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "this is password",
					Phone:    "15212345678",
					Ctime:    now,
				}).Return(nil)

				return ud, uc
			},
			ctx: context.Background(),
			id:  int64(123),

			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "this is password",
				Phone:    "15212345678",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "success, cache hit",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "this is password",
					Phone:    "15212345678",
					Ctime:    now,
				}, nil)
				return nil, uc
			},
			ctx: context.Background(),
			id:  int64(123),

			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "this is password",
				Phone:    "15212345678",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "cache miss, failed to query db",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := daomocks.NewMockUserDAO(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.Users{}, errors.New("db error"))
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				return ud, uc
			},
			ctx: context.Background(),
			id:  int64(123),

			wantUser: domain.User{},
			wantErr:  errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ud, uc := tc.mock(ctrl)
			ur := NewCachedUserRepository(ud, uc)

			u, err := ur.FindById(tc.ctx, tc.id)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)

			// 為了測 go routine的 cache.set 因為他是mock調用絕對在1秒內完成，只有sleep別無他法了！
			time.Sleep(time.Second)
		})
	}
}

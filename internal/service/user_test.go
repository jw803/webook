package service

import (
	"context"
	"errors"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository"
	repomocks "github.com/jw803/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_UserService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		//輸入
		email    string
		password string

		//輸出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登入成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$4fRcRdTJdWbTJSYjJftxD.H7rKF76D.AvkDgv2JJWyP4lLmT1VXEe",
						Phone:    "15212345678",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "abc123321",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$4fRcRdTJdWbTJSYjJftxD.H7rKF76D.AvkDgv2JJWyP4lLmT1VXEe",
				Phone:    "15212345678",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用戶不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "abc123321",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "abc123321",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$MN9ZKKIbjLZDyEpCYW19auY7mvOG9pcpiIcUUoZZI6pA6OmKZKOVi",
						Phone:    "15212345678",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "abc123321",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("abc123321"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}

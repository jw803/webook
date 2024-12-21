package service

import (
	"context"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/internal/repository"
	repomocks "github.com/jw803/webook/internal/repository/mocks"
	"github.com/jw803/webook/internal/test/test_ioc"
	"github.com/jw803/webook/internal/test/test_model"
	"github.com/jw803/webook/pkg/errorx"
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
			name: "login success",
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
			name: "user not found",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errorx.WithCode(test_model.ErrUserNotFound, ""))
				return repo
			},
			email:    "123@qq.com",
			password: "abc123321",

			wantUser: domain.User{},
			wantErr:  errorx.WithCode(test_model.ErrInvalidUserNameOrPassword, ""),
		},
		{
			name: "DB Error",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errorx.WithCode(test_model.ErrDatabase, ""))
				return repo
			},
			email:    "123@qq.com",
			password: "abc123321",

			wantUser: domain.User{},
			wantErr:  errorx.WithCode(test_model.ErrDatabase, ""),
		},
		{
			name: "incorrect password",
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
			wantErr:  errorx.WithCode(errcode.ErrInvalidUserNameOrPassword, "the password user inputted is incorrect"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tc.mock(ctrl), test_ioc.InitLog())
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.True(t, errorx.IsEqual(tc.wantErr, err))
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

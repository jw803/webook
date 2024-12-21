package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/service"
	svcmocks "github.com/jw803/webook/internal/service/mocks"
	"github.com/jw803/webook/internal/test/test_model"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (s *userHandlerSuite) TestUserSignUp() {
	t := s.T()

	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantRes  test_model.Response
	}{
		{
			name: "sign up successfully",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@gmail.com",
					Password: "hello#world123",
				}).Return(nil)
				return userSvc
			},
			reqBody: `
				{
					"email": "123@gmail.com",
					"password": "hello#world123",
					"confirmPassword": "hello#world123"
				}
			`,
			wantCode: http.StatusOK,
			wantRes: test_model.Response{
				Code: test_model.ErrSuccess,
				Data: nil,
				Msg:  "",
			},
		},
		{
			name: "The parameters are incorrect, binding failed",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
				{
					"email": "123@gmail.com",
					"password": "hello#world123",
					"confirmPassword": 
				}
			`,
			wantCode: http.StatusBadRequest,
			wantRes: test_model.Response{
				Code: test_model.ErrBind,
				Data: nil,
				Msg:  "Invalid request format",
			},
		},
		{
			name: "invalid email format",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
				{
					"email": "123@g",
					"password": "hello#world123",
					"confirmPassword": "hello#world123"
				}
			`,
			wantCode: http.StatusBadRequest,
			wantRes: test_model.Response{
				Code: test_model.ErrValidation,
				Data: nil,
				Msg:  "Validation failed",
			},
		},
		{
			name: "password and the confirmation password do not match",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
				{
					"email": "123@gmail.com",
					"password": "hello#world123",
					"confirmPassword": "hello#world"
				}
			`,
			wantCode: http.StatusBadRequest,
			wantRes: test_model.Response{
				Code: test_model.ErrPasswordNotMatch,
				Data: nil,
				Msg:  "The password and the confirmation password do not match",
			},
		},
		{
			name: "invalid password format",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
				{
					"email": "123@gmail.com",
					"password": "hello",
					"confirmPassword": "hello"
				}
			`,
			wantCode: http.StatusBadRequest,
			wantRes: test_model.Response{
				Code: test_model.ErrInvalidPassword,
				Data: nil,
				Msg:  "Invalid password format",
			},
		},
		{
			name: "use existing email to sign up",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@gmail.com",
					Password: "hello#world123",
				}).Return(errorx.WithCode(test_model.ErrDuplicateEmailSignUp, ""))
				return userSvc
			},
			reqBody: `
				{
					"email": "123@gmail.com",
					"password": "hello#world123",
					"confirmPassword": "hello#world123"
				}
			`,
			wantCode: http.StatusBadRequest,
			wantRes: test_model.Response{
				Code: test_model.ErrDuplicateEmailSignUp,
				Data: nil,
				Msg:  "email has already been registered",
			},
		},
		{
			name: "user svc occurs error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@gmail.com",
					Password: "hello#world123",
				}).Return(errors.New("system error"))
				return userSvc
			},
			reqBody: `
				{
					"email": "123@gmail.com",
					"password": "hello#world123",
					"confirmPassword": "hello#world123"
				}
			`,
			wantCode: http.StatusInternalServerError,
			wantRes: test_model.Response{
				Code: test_model.ErrSystem,
				Data: nil,
				Msg:  "Internal server error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer func() {
				ctrl.Finish()
				r := recover()
				if r != nil {
					t.Log("Panic arose! Recovered in ", r)
					t.Fail()
				}
			}()

			gin.SetMode(gin.ReleaseMode)
			server := gin.Default()

			userServiceMock := tc.mock(ctrl)
			h := NewUserHandler(userServiceMock, nil, nil, s.logger)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			var result test_model.Response
			err = json.Unmarshal(resp.Body.Bytes(), &result)
			if err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantRes, result)
		})
	}
}

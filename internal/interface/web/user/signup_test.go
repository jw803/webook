package user

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/service"
	svcmocks "github.com/jw803/webook/internal/service/mocks"
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
		wantBody string
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
			wantBody: "{}",
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

			server := gin.Default()

			userServiceMock := tc.mock(ctrl)
			h := NewUserHandler(userServiceMock, nil, nil, s.logger)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

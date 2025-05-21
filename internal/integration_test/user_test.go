package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/config"
	"github.com/jw803/webook/internal/integration_test/startup"
	"github.com/jw803/webook/internal/interface/web"
	"github.com/jw803/webook/internal/pkg/ginx"
	"github.com/jw803/webook/internal/repository/dao"
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type WebAPIUserTestSuite struct {
	suite.Suite
	db     *gorm.DB
	redis  redis.Cmdable
	logger loggerx.Logger
}

func TestWebAPIUser(t *testing.T) {
	suite.Run(t, new(WebAPIUserTestSuite))
}

func (s *WebAPIUserTestSuite) SetupSuite() {
	config.SetTestConfig()
	s.logger = startup.InitLogger()
	s.db = startup.InitDB()
	s.redis = startup.InitRedis()
}

func (s *WebAPIUserTestSuite) setServer(handler web.Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	handler.RegisterRoutes(server)
	return server
}

func (s *WebAPIUserTestSuite) ClearTestCaseData() {
	err := s.db.Exec("truncate table `articles`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("truncate table `publish_articles`").Error
	assert.NoError(s.T(), err)
}

func (s *WebAPIUserTestSuite) TestUserHandler_SendSMSCode() {
	t := s.T()
	testCases := []struct {
		name   string
		mock   func(t *testing.T, ctrl *gomock.Controller) *gin.Engine
		before func(t *testing.T)
		after  func(t *testing.T)

		phone string

		wantCode int
		wantBody ginx.Response
	}{
		{
			name: "发送成功的用例",
			mock: func(t *testing.T, ctrl *gomock.Controller) *gin.Engine {
				nowFunc := startup.NewNowFunc("2025-01-01T00:00:00Z")
				userDao := dao.NewGORMUserDAO(s.db, s.logger, nowFunc)
				userHandler := startup.InitUserhandler(userDao)
				return s.setServer(userHandler)
				server := gin.New()
				return server
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := s.redis.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := s.redis.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				err = s.redis.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: ginx.Response{
				Message: "",
			},
		},
		{
			name: "未输入手机号码",
			before: func(t *testing.T) {

			},
			after:    func(t *testing.T) {},
			wantCode: http.StatusOK,
			wantBody: ginx.Response{
				Code:    4,
				Message: "请输入手机号码",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				err := s.redis.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := s.redis.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: ginx.Response{
				Code:    4,
				Message: "短信发送太频繁，请稍后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				err := s.redis.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := s.redis.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: ginx.Response{
				Code:    5,
				Message: "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer func() {
				r := recover()
				if r != nil {
					t.Log("Recovered in ", r)
					t.Fail()
				}
				s.ClearTestCaseData()
			}()
			defer tc.after(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := tc.mock(t, ctrl)
			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone": "%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res ginx.Response
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}

}

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/jw803/webook/internal/integration/startup"
	"github.com/jw803/webook/internal/web"
	"github.com/jw803/webook/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := startup.InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		before  func(t *testing.T)
		after   func(t *testing.T)
		reqBody string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "發送成功",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:12345678").Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody: `
				{
					"phone": "12345678"
				}
			`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "發送太頻繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:12345678", "123456", time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:12345678").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
			},
			reqBody: `
				{
					"phone": "12345678"
				}
			`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送太频繁，请稍后再试",
			},
		},
		{
			name: "系統錯誤",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:12345678", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:12345678").Result()
				cancel()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
			},
			reqBody: `
				{
					"phone": "12345678"
				}
			`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手機號碼為空",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
				{
					"phone": ""
				}
			`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 4,
				Msg:  "输入有误",
			},
		},
		{
			name: "req數據格式錯誤",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
				{
					"phone":
				}
			`,
			wantCode: 400,
			wantBody: web.Result{
				Code: 4,
				Msg:  "输入有误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			if resp.Code != 200 {
				return
			}
			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)

			assert.Equal(t, tc.wantBody, webRes)
			tc.after(t)
		})
	}
}

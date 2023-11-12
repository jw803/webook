package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/service"
	svcmocks "github.com/jw803/webook/internal/service/mocks"
	ijwt "github.com/jw803/webook/internal/web/jwt"
	"github.com/jw803/webook/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "新建並發表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return articleSvc
			},
			reqBody: `
				{
					"title": "我的標題",
					"content": "我的內容"
				}
			`,
			wantCode: 200,
			wantRes: Result{
				Data: float64(1),
				Msg:  "OK",
			},
		},
		{
			name: "publish失敗",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish 失敗"))
				return articleSvc
			},
			reqBody: `
				{
					"title": "我的標題",
					"content": "我的內容"
				}
			`,
			wantCode: 200,
			wantRes: Result{
				Code: 5,
				Msg:  "系統錯誤",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("claims", &ijwt.UserClaims{
					Uid: 123,
				})
			})

			h := NewArticleHandler(tc.mock(ctrl), logger.NewNoOpLogger())
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			if resp.Code != 200 {
				return
			}
			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)

			assert.Equal(t, tc.wantRes, webRes)
			assert.Equal(t, tc.wantCode, resp.Code)
		})
	}
}

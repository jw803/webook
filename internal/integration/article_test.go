package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/repository/dao"
	ijwt "github.com/jw803/webook/internal/web/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jw803/webook/internal/integration/startup"
)

type ArticleTestSuit struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuit) SetupSuite() {
	// 全部路由都註冊
	//s.server = startup.InitWebServer()
	s.db = startup.InitTestDB()
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			Uid: 123,
		})
	})
	// 這個還要處理service注入log的問題
	//articleHdl := web.NewArticleHandler(service.NewArticleService())
	articleHdl := startup.InitArticleHandler()
	articleHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuit) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE articles")
}
func (s *ArticleTestSuit) TestArticle() {
	t := s.T()
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		article Article

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建帖子 - 保存成功",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.db.Where("Id=?", 1).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)

				article.Ctime = 0
				article.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "我的標題",
					Content:  "我的內容",
					AuthorId: 123,
				}, article)
			},
			article: Article{
				Title:   "我的標題",
				Content: "我的內容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			// 這樣做會測不到非法json，但其實gin會幫我們處理所以不測也ok
			reqBody, err := json.Marshal(tc.article)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			if resp.Code != 200 {
				return
			}
			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)

			assert.Equal(t, tc.wantRes, webRes)
			tc.after(t)
		})
	}
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func (s *ArticleTestSuit) TestABC() {
	s.T().Log("Hello 這是測試套件")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuit{})
}

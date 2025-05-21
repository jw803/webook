package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	articleDao "github.com/jw803/webook/internal/repository/dao/article"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jw803/webook/internal/test/startup"
)

type ArticleGORMHandlerTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleGORMHandlerTestSuite) SetupSuite() {
	// 全部路由都註冊
	//s.server = startup.InitWebServer()
	s.db = startup.InitTestDB()
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &jwtx.UserClaims{
			Uid: 123,
		})
		ctx.Next()
	})
	// 這個還要處理service注入log的問題
	//articleHdl := web.NewArticleHandler(service.NewArticleService())
	s.db = startup.InitTestDB()
	articleHdl := startup.InitArticleHandler(articleDao.NewGORMArticleDAO(s.db))
	articleHdl.RegisterRoutes(s.server)
}

func (s *ArticleGORMHandlerTestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE articles")
	s.db.Exec("TRUNCATE TABLE published_articles")
}

func (s *ArticleGORMHandlerTestSuite) TestArticleHandler_Edit() {
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
				var article articleDao.Article
				err := s.db.Where("Id=?", 1).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)

				article.Ctime = 0
				article.Utime = 0
				assert.Equal(t, articleDao.Article{
					Id:       1,
					Title:    "我的標題",
					Content:  "我的內容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
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
		{
			name: "修改已有帖子，並保存",
			before: func(t *testing.T) {
				s.db.Create(&articleDao.Article{
					Id:       2,
					Title:    "我的標題",
					Content:  "我的內容",
					Ctime:    111,
					Utime:    111,
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				})
			},
			after: func(t *testing.T) {
				var article articleDao.Article
				err := s.db.Where("Id=?", 2).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Utime > 111)

				article.Utime = 0
				assert.Equal(t, articleDao.Article{
					Id:       2,
					Title:    "新的標題",
					Content:  "新的內容",
					Ctime:    111,
					Utime:    0,
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}, article)
			},
			article: Article{
				Id:      2,
				Title:   "新的標題",
				Content: "新的內容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
				Msg:  "OK",
			},
		},
		{
			name: "非文章作者本人跑來修改文章",
			before: func(t *testing.T) {
				s.db.Create(&articleDao.Article{
					Id:       3,
					Title:    "我的標題",
					Content:  "我的內容",
					Ctime:    111,
					Utime:    111,
					AuthorId: 456,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				})
			},
			after: func(t *testing.T) {
				var article articleDao.Article
				err := s.db.Where("Id=?", 3).First(&article).Error
				assert.NoError(t, err)
				assert.Equal(t, articleDao.Article{
					Id:       3,
					Title:    "我的標題",
					Content:  "我的內容",
					Ctime:    111,
					Utime:    111,
					AuthorId: 456,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}, article)
			},
			article: Article{
				Id:      3,
				Title:   "新的標題",
				Content: "新的內容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系統錯誤",
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

			if resp.Code != 200 {
				return
			}
			assert.Equal(t, tc.wantCode, resp.Code)

			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)

			assert.Equal(t, tc.wantRes, webRes)
			tc.after(t)
		})
	}
}

func (s *ArticleGORMHandlerTestSuite) TestArticleHandler_Publish() {
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
			name: "新建帖子並發表",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				var article articleDao.Article
				err := s.db.Where("Id=?", 1).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)
				article.Ctime = 0
				article.Utime = 0
				assert.Equal(t, articleDao.Article{
					Id:       1,
					Title:    "我的標題",
					Content:  "我的內容",
					AuthorId: 123,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}, article)

				var readerArticle articleDao.PublishArticle
				err = s.db.Where("Id=?", 1).First(&readerArticle).Error
				assert.NoError(t, err)
				assert.True(t, readerArticle.Ctime > 0)
				assert.True(t, readerArticle.Utime > 0)
				readerArticle.Ctime = 0
				readerArticle.Utime = 0
				assert.Equal(t, articleDao.PublishArticle{
					Id:       1,
					Title:    "我的標題",
					Content:  "我的內容",
					AuthorId: 123,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}, readerArticle)
			},
			article: Article{
				Title:   "我的標題",
				Content: "我的內容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "OK",
				Data: 1,
			},
		},
		{
			name: "更新帖子並發表",
			before: func(t *testing.T) {
				s.db.Create(articleDao.Article{
					Id:       2,
					Title:    "我的標題",
					Content:  "我的內容",
					Ctime:    111,
					Utime:    111,
					AuthorId: 123,
					Status:   articleDao.ArticleStatusUnpublished.ToUint8(),
				})
			},
			after: func(t *testing.T) {
				var article articleDao.Article
				err := s.db.Where("id = ?", 2).First(&article).Error
				assert.NoError(t, err)

				assert.True(t, article.Utime > 111)
				article.Utime = 0
				assert.Equal(t, articleDao.Article{
					Id:       2,
					Title:    "新的標題",
					Content:  "新的內容",
					AuthorId: 123,
					Ctime:    111,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}, article)

				var publishArticle articleDao.PublishArticle
				err = s.db.Where("id = ?", 2).First(&publishArticle).Error
				assert.NoError(t, err)

				assert.True(t, publishArticle.Ctime > 0)
				assert.True(t, publishArticle.Utime > 0)
				publishArticle.Ctime = 0
				publishArticle.Utime = 0
				assert.Equal(t, articleDao.PublishArticle{
					Id:       2,
					Title:    "新的標題",
					Content:  "新的內容",
					AuthorId: 123,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}, publishArticle)

			},

			article: Article{
				Id:      2,
				Title:   "新的標題",
				Content: "新的內容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "OK",
				Data: 2,
			},
		},
		{
			name: "更新帖子並重新發表",
			before: func(t *testing.T) {
				article := articleDao.Article{
					Id:       3,
					Title:    "我的標題",
					Content:  "我的內容",
					Ctime:    111,
					Utime:    111,
					AuthorId: 123,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}
				s.db.Create(&article)
				publishArticle := articleDao.PublishArticle{
					Id:       3,
					Title:    "我的標題",
					Content:  "我的內容",
					Ctime:    111,
					Utime:    111,
					AuthorId: 123,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}
				s.db.Create(&publishArticle)
			},
			after: func(t *testing.T) {
				var article articleDao.Article
				err := s.db.Where("id = ?", 3).First(&article).Error
				assert.NoError(t, err)

				assert.True(t, article.Utime > 111)
				article.Utime = 0
				assert.Equal(t, articleDao.Article{
					Id:       3,
					Title:    "新的標題",
					Content:  "新的內容",
					AuthorId: 123,
					Ctime:    111,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}, article)

				var publishArticle articleDao.PublishArticle
				err = s.db.Where("id = ?", 3).First(&publishArticle).Error
				assert.NoError(t, err)

				assert.True(t, publishArticle.Utime > 111)
				publishArticle.Utime = 0

				assert.Equal(t, articleDao.PublishArticle{
					Id:       3,
					Title:    "新的標題",
					Content:  "新的內容",
					AuthorId: 123,
					Ctime:    111,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}, publishArticle)
			},

			article: Article{
				Id:      3,
				Title:   "新的標題",
				Content: "新的內容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "OK",
				Data: 3,
			},
		},
		{
			name: "更新别人的帖子，并且发表失败",
			before: func(t *testing.T) {
				article := articleDao.Article{
					Id:       4,
					Title:    "我的標題",
					Content:  "我的內容",
					Ctime:    111,
					Utime:    111,
					AuthorId: 555,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}
				s.db.Create(&article)
				publishArticle := articleDao.PublishArticle{
					Id:       4,
					Title:    "我的標題",
					Content:  "我的內容",
					Ctime:    111,
					Utime:    111,
					AuthorId: 555,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}
				s.db.Create(&publishArticle)
			},
			after: func(t *testing.T) {
				var article articleDao.Article
				err := s.db.Where("id = ?", 4).First(&article).Error
				assert.NoError(t, err)

				assert.Equal(t, articleDao.Article{
					Id:       4,
					Title:    "我的標題",
					Content:  "我的內容",
					AuthorId: 555,
					Ctime:    111,
					Utime:    111,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}, article)

				var publishArticle articleDao.PublishArticle
				err = s.db.Where("id = ?", 4).First(&publishArticle).Error
				assert.NoError(t, err)

				assert.Equal(t, articleDao.PublishArticle{
					Id:       4,
					Title:    "我的標題",
					Content:  "我的內容",
					AuthorId: 555,
					Ctime:    111,
					Utime:    111,
					Status:   articleDao.ArticleStatusPublished.ToUint8(),
				}, publishArticle)
			},

			article: Article{
				Id:      4,
				Title:   "新的標題",
				Content: "新的內容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系統錯誤",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			// 這樣做會測不到非法json，但其實gin會幫我們處理所以不測也ok
			reqBody, err := json.Marshal(tc.article)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer(reqBody))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				return
			}
			assert.Equal(t, tc.wantCode, resp.Code)

			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)

			tc.after(t)
		})
	}
}

func (s *ArticleGORMHandlerTestSuite) TestABC() {
	s.T().Log("Hello 這是測試套件")
}

func TestGORMArticle(t *testing.T) {
	suite.Run(t, &ArticleGORMHandlerTestSuite{})
}

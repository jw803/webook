package service

import (
	"context"
	"errors"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository/article"
	repoarticlemocks "github.com/jw803/webook/internal/repository/article/mocks"
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository)

		ctx context.Context

		article domain.Article

		wantId  int64
		wantErr error
	}{
		{
			name: "新建並發表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repoarticlemocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				readerRepo := repoarticlemocks.NewMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return authorRepo, readerRepo
			},
			ctx: context.Background(),
			article: domain.Article{
				Title:   "我的標題",
				Content: "我的內容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "修改並發表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repoarticlemocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepo := repoarticlemocks.NewMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)
				return authorRepo, readerRepo
			},
			ctx: context.Background(),
			article: domain.Article{
				Id:      2,
				Title:   "我的標題",
				Content: "我的內容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  2,
			wantErr: nil,
		},
		{
			name: "保存到製作庫失敗",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repoarticlemocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("mock db error"))
				return authorRepo, nil
			},
			ctx: context.Background(),
			article: domain.Article{
				Id:      2,
				Title:   "我的標題",
				Content: "我的內容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("mock db error"),
		},
		{
			name: "保存到製作庫成功，且經過重試，線上庫也成功了",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repoarticlemocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepo := repoarticlemocks.NewMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("mock db error"))
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)
				return authorRepo, readerRepo
			},
			ctx: context.Background(),
			article: domain.Article{
				Id:      2,
				Title:   "我的標題",
				Content: "我的內容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  2,
			wantErr: nil,
		},
		{
			name: "保存到製作庫成功，且經過重試，最終線上庫全部重試失敗",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repoarticlemocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepo := repoarticlemocks.NewMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的標題",
					Content: "我的內容",
					Author: domain.Author{
						Id: 123,
					},
				}).Times(3).Return(int64(0), errors.New("mock db error"))
				return authorRepo, readerRepo
			},
			ctx: context.Background(),
			article: domain.Article{
				Id:      2,
				Title:   "我的標題",
				Content: "我的內容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authorRepo, readerRepo := tc.mock(ctrl)
			svc := NewArticleServiceV1(authorRepo, readerRepo, loggerx.NewNoOpLogger())
			id, err := svc.PublishV1(tc.ctx, tc.article)

			assert.Equal(t, tc.wantId, id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

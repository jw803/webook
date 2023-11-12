package service

import (
	"context"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository/article"
	"github.com/jw803/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	PublishV1(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo article.ArticleRepository

	// v1
	authorRepo article.ArticleAuthorRepository
	readerRepo article.ArticleReaderRepository
	l          logger.LoggerV1
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(authorRepo article.ArticleAuthorRepository, readerRepo article.ArticleReaderRepository,
	l logger.LoggerV1) ArticleService {
	return &articleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		l:          l,
	}
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	} else {
		return a.repo.Create(ctx, article)
	}
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	return a.repo.SyncV1(ctx, article)
}

func (a *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	if article.Id > 0 {
		err = a.authorRepo.Update(ctx, article)
	} else {
		id, err = a.authorRepo.Create(ctx, article)
	}

	if err != nil {
		return 0, err
	}
	article.Id = id
	for i := 0; i < 3; i++ {
		id, err = a.readerRepo.Save(ctx, article)
		if err == nil {
			break
		}
		a.l.Error("部分失敗，保存到線上庫失敗", logger.Int64("art_id", article.Id),
			logger.Error(err))
	}
	if err != nil {
		a.l.Error("部分失敗，保存到線上庫的重試全部失敗", logger.Int64("art_id", article.Id),
			logger.Error(err))
		return 0, err
	}
	return id, err
}

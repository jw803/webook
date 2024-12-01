package service

import (
	"context"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository/article"
	"github.com/jw803/webook/pkg/loggerx"
	"time"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	PublishV1(ctx context.Context, article domain.Article) (int64, error)

	Withdraw(ctx context.Context, art domain.Article) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	// ListPub 根据这个 start 时间来查询
	ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error)
}

type articleService struct {
	repo article.ArticleRepository

	// v1
	authorRepo article.ArticleAuthorRepository
	readerRepo article.ArticleReaderRepository
	l          loggerx.Logger
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(authorRepo article.ArticleAuthorRepository, readerRepo article.ArticleReaderRepository,
	l loggerx.Logger) ArticleService {
	return &articleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		l:          l,
	}
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnpublished
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	} else {
		return a.repo.Create(ctx, article)
	}
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, article)
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
		a.l.Error(ctx, "部分失敗，保存到線上庫失敗", loggerx.Int64("art_id", article.Id),
			loggerx.Error(err))
	}
	if err != nil {
		a.l.Error(ctx, "部分失敗，保存到線上庫的重試全部失敗", loggerx.Int64("art_id", article.Id),
			loggerx.Error(err))
		return 0, err
	}
	return id, err
}

func (a *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	// art.Status = domain.ArticleStatusPrivate 然后直接把整个 art 往下传
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (svc *articleService) ListPub(ctx context.Context,
	start time.Time, offset, limit int) ([]domain.Article, error) {
	return svc.repo.ListPub(ctx, start, offset, limit)
}

func (svc *articleService) GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error) {
	// 另一个选项，在这里组装 Author，调用 UserService
	art, err := svc.repo.GetPublishedById(ctx, id)
	return art, err
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *articleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

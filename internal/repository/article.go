package repository

import (
	"context"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDao
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

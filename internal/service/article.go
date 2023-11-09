package service

import (
	"context"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func (a articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	return a.repo.Create(ctx, article)
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

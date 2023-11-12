package article

import (
	"context"
	"github.com/jw803/webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}

package article

import (
	"context"
	"gorm.io/gorm"
)

type ReaderDao interface {
	Upsert(ctx context.Context, article Article) error
}

func NewReaderDao(db *gorm.DB) ReaderDao {
	return nil
}

// 這個代表的是線上表
type PublishArticle struct {
	Article
}

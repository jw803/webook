package article

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleReaderDAO interface {
	Upsert(ctx context.Context, article PublishArticle) error
}

type GORMArticleReaderDAO struct {
	db *gorm.DB
}

func NewReaderDao(db *gorm.DB) ArticleReaderDAO {
	return &GORMArticleReaderDAO{
		db: db,
	}
}

// 這個代表的是線上表
type PublishArticle Article

func (dao *GORMArticleReaderDAO) Upsert(ctx context.Context, article PublishArticle) error {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := dao.db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   article.Utime,
			"status":  article.Status,
		}),
	}).Create(&article).Error
	return err
}

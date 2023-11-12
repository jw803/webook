package article

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Upsert(ctx context.Context, article PublishArticle) error

	Sync(ctx context.Context, article Article) (int64, error)
}

func NewGORMArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleDao{
		db: db,
	}
}

type GORMArticleDao struct {
	db *gorm.DB
}

func (dao *GORMArticleDao) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := dao.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func (dao *GORMArticleDao) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Utime = now
	res := dao.db.WithContext(ctx).Model(&article).
		Where("id=? and author_id=?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   article.Utime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("可能是非作者想來更新這篇文章")
	}
	return nil
}

func (dao *GORMArticleDao) Upsert(ctx context.Context, article PublishArticle) error {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := dao.db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   article.Utime,
		}),
	}).Create(&article).Error
	return err
}

func (dao *GORMArticleDao) Sync(ctx context.Context, article Article) (int64, error) {
	var (
		id = article.Id
	)
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewGORMArticleDao(tx)
		if id > 0 {
			err = txDAO.UpdateById(ctx, article)
		} else {
			id, err = txDAO.Insert(ctx, article)
		}
		if err != nil {
			return err
		}

		return txDAO.Upsert(ctx, PublishArticle{Article: article})
	})
	return id, err
}

type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 标题的长度
	// 正常都不会超过这个长度
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`
	// 作者
	AuthorId int64 `gorm:"index"`
	Ctime    int64
	Utime    int64
}

package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
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

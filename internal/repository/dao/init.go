package dao

import (
	"github.com/jw803/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&Users{}, &article.Article{}, &article.PublishArticle{})
}

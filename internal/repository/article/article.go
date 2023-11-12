package article

import (
	"context"
	"github.com/jw803/webook/internal/domain"
	dao "github.com/jw803/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error

	SyncV1(ctx context.Context, article domain.Article) (int64, error)
	SyncV2(ctx context.Context, article domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDao

	//V1
	authorDao dao.AuthorDao
	readerDao dao.ReaderDao

	//V2 缺點是誒何了dao操作的東西
	db *gorm.DB
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func (c *CachedArticleRepository) SyncV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	articlen := c.ToEntity(article)
	if id > 0 {
		err = c.authorDao.UpdateById(ctx, articlen)
	} else {
		id, err = c.authorDao.Insert(ctx, articlen)
	}
	if err != nil {
		return id, err
	}

	err = c.readerDao.Upsert(ctx, articlen)
	return id, err
}

func (c *CachedArticleRepository) SyncV2(ctx context.Context, article domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 怕中間有panic 會不小心沒rollback，因為err判斷是處理不到的
	defer tx.Rollback()
	author := dao.NewAuthorDao(tx)
	reader := dao.NewReaderDao(tx)

	var (
		id  = article.Id
		err error
	)
	articlen := c.ToEntity(article)
	if id > 0 {
		err = author.UpdateById(ctx, articlen)
	} else {
		id, err = author.Insert(ctx, articlen)
	}
	if err != nil {
		//有defer就不用了
		//tx.Rollback()
		return id, err
	}

	err = reader.Upsert(ctx, articlen)
	tx.Commit()
	return id, err
}

func (c *CachedArticleRepository) ToEntity(article domain.Article) dao.Article {
	return dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	}
}

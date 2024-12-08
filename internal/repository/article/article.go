package article

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/internal/repository/cache"
	articleDao "github.com/jw803/webook/internal/repository/dao/article"
	"github.com/jw803/webook/pkg/loggerx"
	"gorm.io/gorm"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error

	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncV1(ctx context.Context, article domain.Article) (int64, error)
	SyncV2(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, authorId int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	userRepo repository.UserRepository
	dao      articleDao.ArticleDAO

	//V1
	authorDao articleDao.AuthorDao
	readerDao articleDao.ArticleReaderDAO

	//V2 缺點是誒何了dao操作的東西
	db *gorm.DB

	cache cache.ArticleCache
	l     loggerx.Logger
}

func NewArticleRepository(dao articleDao.ArticleDAO, articleCache cache.ArticleCache) ArticleRepository {
	return &CachedArticleRepository{
		dao:   dao,
		cache: articleCache,
	}
}

func (repo *CachedArticleRepository) ToEntity(article domain.Article) articleDao.Article {
	return articleDao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	}
}

func (repo *CachedArticleRepository) toDomain(art articleDao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}

func (repo *CachedArticleRepository) preCache(ctx context.Context, data []domain.Article) {
	if len(data) > 0 && len(data[0].Content) < 1024*1024 {
		err := repo.cache.Set(ctx, data[0])
		if err != nil {
			repo.l.Error(ctx, "提前预加载缓存失败", loggerx.Error(err))
		}
	}
}

func (repo *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	defer func() {
		// 清空缓存
		repo.cache.DelFirstPage(ctx, article.Author.Id)
	}()
	return repo.dao.Insert(ctx, repo.ToEntity(article))
}

func (repo *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	defer func() {
		// 清空缓存
		repo.cache.DelFirstPage(ctx, article.Author.Id)
	}()
	return repo.dao.UpdateById(ctx, repo.ToEntity(article))
}

func (repo *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := repo.dao.Sync(ctx, repo.ToEntity(art))
	if err == nil {
		repo.cache.DelFirstPage(ctx, art.Author.Id)
		er := repo.cache.SetPub(ctx, art)
		if er != nil {
			// 不需要特别关心
			// 比如说输出 WARN 日志
			repo.l.Warn(ctx, "failed to set pub article", loggerx.Error(er))
		}
	}
	return id, err
}

func (repo *CachedArticleRepository) SyncV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	articlen := repo.ToEntity(article)
	if id > 0 {
		err = repo.authorDao.UpdateById(ctx, articlen)
	} else {
		id, err = repo.authorDao.Insert(ctx, articlen)
	}
	if err != nil {
		return id, err
	}

	err = repo.readerDao.Upsert(ctx, articleDao.PublishArticle(articlen))
	return id, err
}

func (repo *CachedArticleRepository) SyncV2(ctx context.Context, article domain.Article) (int64, error) {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 怕中間有panic 會不小心沒rollback，因為err判斷是處理不到的
	defer tx.Rollback()
	author := articleDao.NewAuthorDao(tx)
	reader := articleDao.NewReaderDao(tx)

	var (
		id  = article.Id
		err error
	)
	articlen := repo.ToEntity(article)
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

	err = reader.Upsert(ctx, articleDao.PublishArticle(articlen))
	tx.Commit()
	return id, err
}

func (repo *CachedArticleRepository) SyncStatus(ctx context.Context, id int64, authorId int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, id, authorId, uint8(status))
}

func (repo *CachedArticleRepository) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Article, error) {
	res, err := repo.dao.ListPub(ctx, start, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(res, func(idx int, src articleDao.Article) domain.Article {
		return repo.toDomain(src)
	}), nil
}

func (repo *CachedArticleRepository) GetPublishedById(
	ctx context.Context, id int64) (domain.Article, error) {
	// 读取线上库数据，如果你的 Content 被你放过去了 OSS 上，你就要让前端去读 Content 字段
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 你在这边要组装 user 了，适合单体应用
	usr, err := repo.userRepo.FindById(ctx, art.AuthorId)
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.NickName,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	return res, nil
}

func (repo *CachedArticleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	data, err := repo.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return repo.toDomain(data), nil
}

func (repo *CachedArticleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 你在这个地方，集成你的复杂的缓存方案
	// 你只缓存这一页
	if offset == 0 && limit <= 100 {
		data, err := repo.cache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				repo.preCache(ctx, data)
			}()
			//return data[:limit], err
			return data, err
		}
	}
	res, err := repo.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	data := slice.Map[articleDao.Article, domain.Article](res, func(idx int, src articleDao.Article) domain.Article {
		return repo.toDomain(src)
	})
	// 回写缓存的时候，可以同步，也可以异步
	go func() {
		err := repo.cache.SetFirstPage(ctx, uid, data)
		repo.l.Error(ctx, "回写缓存失败", loggerx.Error(err))
		repo.preCache(ctx, data)
	}()
	return data, nil
}

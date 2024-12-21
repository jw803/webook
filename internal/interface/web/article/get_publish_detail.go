package article

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"golang.org/x/sync/errgroup"
	"time"
)

type getPubArticleParam struct {
	Id int64 `uri:"id"`
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context, uc *jwtx.UserClaims, params getPubArticleParam) (any, error) {
	var article articleVO
	id := params.Id

	var eg errgroup.Group
	var art domain.Article
	var err error
	eg.Go(func() error {
		art, err = h.svc.GetPublishedById(ctx, id, uc.Uid)
		return err
	})

	// 在这儿等，要保证前面两个
	err = eg.Wait()
	if err != nil {
		h.l.P1(ctx, "failed to get pub article detail")
		return "", err
	}

	// ctx.Set("art", art)

	article = articleVO{
		Id:      art.Id,
		Title:   art.Title,
		Status:  art.Status.ToUint8(),
		Content: art.Content,
		// 要把作者信息带出去
		Author: art.Author.Name,
		Ctime:  art.Ctime.Format(time.DateTime),
		Utime:  art.Utime.Format(time.DateTime),
	}
	// 这个功能是不是可以让前端，主动发一个 HTTP 请求，来增加一个计数？
	return article, nil
}

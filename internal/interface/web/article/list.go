package article

import (
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"time"
)

type listQuery struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (h *ArticleHandler) List(ctx *gin.Context, uc *jwtx.UserClaims, req listQuery) (any, error) {
	result := make([]articleVO, 0)
	res, err := h.svc.List(ctx, uc.Uid, req.Offset, req.Limit)
	if err != nil {
		h.l.P1(ctx, "failed to edit article")
		return result, err
	}

	// 在列表页，不显示全文，只显示一个"摘要"
	// 比如说，简单的摘要就是前几句话
	// 强大的摘要是 AI 帮你生成的
	result = slice.Map[domain.Article, articleVO](res,
		func(idx int, src domain.Article) articleVO {
			return articleVO{
				Id:       src.Id,
				Title:    src.Title,
				Abstract: src.Abstract(),
				Status:   src.Status.ToUint8(),
				// 这个列表请求，不需要返回内容
				//Content: src.Content,
				// 这个是创作者看自己的文章列表，也不需要这个字段
				//Author: src.Author
				Ctime: src.Ctime.Format(time.DateTime),
				Utime: src.Utime.Format(time.DateTime),
			}
		})
	return result, nil
}

package article

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	ijwt "github.com/jw803/webook/internal/interface/web/jwtx"
)

type withdrawArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context, uc *ijwt.UserClaims, req withdrawArticleReq) (any, error) {
	err := h.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		h.l.P1(ctx, "failed to withdraw article")
		return nil, err
	}
	return nil, nil
}
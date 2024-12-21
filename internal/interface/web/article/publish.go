package article

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
)

type publishArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req publishArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}

func (h *ArticleHandler) Publish(ctx *gin.Context, uc *jwtx.UserClaims, req publishArticleReq) (any, error) {
	id, err := h.svc.Publish(ctx, req.toDomain(uc.Uid))
	if err != nil {
		h.l.P1(ctx, "failed to edit article")
		return "", err
	}
	return id, nil
}

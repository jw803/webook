package article

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/pkg/errcode"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/pkg/errorx"
)

type editArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req editArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}

func (h *ArticleHandler) Edit(ctx *gin.Context, uc *jwtx.UserClaims, req editArticleReq) (any, error) {
	id, err := h.svc.Save(ctx, req.toDomain(uc.Uid))
	if errorx.IsCode(err, errcode.ErrArticleNotFound) {
		h.l.P2(ctx, "failed to edit article")
		return "", errorx.WithCode(errcode.ErrMaliciousUser, "malicious user intend to edit article")
	}
	if err != nil {
		h.l.P1(ctx, "failed to edit article")
		return "", err
	}
	return id, nil
}

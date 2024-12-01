package article

import (
	"github.com/gin-gonic/gin"
	ijwt "github.com/jw803/webook/internal/interface/web/jwtx"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"time"
)

type getArticleParam struct {
	Id int64 `uri:"id"`
}

func (h *ArticleHandler) Detail(ctx *gin.Context, usr *ijwt.UserClaims, params getArticleParam) (any, error) {
	var result articleVO
	art, err := h.svc.GetById(ctx, params.Id)
	if err != nil {
		h.l.P1(ctx, "failed to get article detail")
		return nil, err
	}
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != usr.Uid {
		//ctx.JSON(http.StatusOK)
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		//a.l.Error("非法访问文章，创作者 ID 不匹配",
		//	logger.Int64("uid", usr.Id))
		h.l.P2(ctx, "malicious user intend to access article")
		return nil, errorx.WithCode(errcode.ErrMaliciousUser, "malicious user intend to access article")
	}
	result = articleVO{
		Id:    art.Id,
		Title: art.Title,
		// 不需要这个摘要信息
		//Abstract: art.Abstract(),
		Status:  art.Status.ToUint8(),
		Content: art.Content,
		// 这个是创作者看自己的文章列表，也不需要这个字段
		//Author: art.Author
		Ctime: art.Ctime.Format(time.DateTime),
		Utime: art.Utime.Format(time.DateTime),
	},
	return result, nil
}

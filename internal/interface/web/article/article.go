package article

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/interface/web"
	ijwt "github.com/jw803/webook/internal/interface/web/jwtx"
	"github.com/jw803/webook/internal/pkg/ginx"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/pkg/loggerx"
)

var _ web.Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   loggerx.Logger
}

func NewArticleHandler(svc service.ArticleService, l loggerx.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/articles")
	ug.POST("/edit", ginx.WrapClaimsReq[ijwt.UserClaims, editArticleReq](h.Edit))
	ug.POST("/publish", ginx.WrapClaimsReq[ijwt.UserClaims, publishArticleReq](h.Publish))
	ug.POST("/withdraw", ginx.WrapClaimsReq[ijwt.UserClaims, withdrawArticleReq](h.Withdraw))
	ug.POST("/list", ginx.WrapClaimsQuery[ijwt.UserClaims, listQuery](h.List))
	ug.GET("/detail/:id", ginx.WrapClaimsParam[ijwt.UserClaims, getArticleParam](h.Detail))

	pub := ug.Group("/pub")
	pub.GET("/:id", ginx.WrapClaimsParam[ijwt.UserClaims, getPubArticleParam](h.PubDetail)
}

type articleVO struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	// 内容
	Content string `json:"content"`
	// 注意一点，状态这个东西，可以是前端来处理，也可以是后端处理
	// 0 -> unknown -> 未知状态
	// 1 -> 未发表，手机 APP 这种涉及到发版的问题，那么后端来处理
	// 涉及到国际化，也是后端来处理
	Status uint8  `json:"status"`
	Author string `json:"author"`
	// 计数
	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`

	// 我个人有没有收藏，有没有点赞
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`

	Ctime string `json:"ctime"`
	Utime string `json:"utime"`
}

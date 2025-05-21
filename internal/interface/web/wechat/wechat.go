package wechat

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/interface/web/wechat/middlewares"
	"github.com/jw803/webook/internal/pkg/ginx"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/internal/service/oauth2/wechat"
	"github.com/jw803/webook/pkg/loggerx"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtx.Handler
	key             []byte
	stateCookieName string
	l               loggerx.Logger
}

func NewOAuth2WechatHandler(svc wechat.Service,
	hdl jwtx.Handler,
	userSvc service.UserService,
	l loggerx.Logger) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"),
		stateCookieName: "jwt-state",
		Handler:         hdl,
		l:               l,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", ginx.Wrap(o.Auth2URL))
	g.Any("/callback",
		middlewares.NewVerifyStateMiddlewareBuilder(o.l).
			SetCookieStateKey("jwt-state").SetQueryStateKey("state").
			Build(),
		ginx.WrapQuery[CallBackQuery](o.Callback))
}

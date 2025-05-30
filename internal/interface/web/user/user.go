package user

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/interface/web"
	"github.com/jw803/webook/internal/pkg/ginx"
	jwtx "github.com/jw803/webook/internal/pkg/ginx/jwt_handler"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/pkg/loggerx"
)

const biz = "login"

var _ web.Handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	passwordExp *regexp.Regexp
	jwtx.JWTHandler
	l loggerx.Logger
}

func NewUserHandler(svc service.UserService,
	codeSvc service.CodeService, jwtHdl jwtx.JWTHandler, l loggerx.Logger) *UserHandler {
	const (
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		passwordExp: passwordExp,
		codeSvc:     codeSvc,
		JWTHandler:  jwtHdl,
		l:           l,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", ginx.WrapReq[signUpReq](h.SignUp))
	ug.POST("/login", ginx.WrapReq[loginReq](h.LoginJWT))
	ug.POST("/logout", ginx.WrapClaim[jwtx.UserClaims](h.LogoutJWT))

	ug.POST("/login_sms/code/send", ginx.WrapReq[loginSMSSendCodeReq](h.SendLoginSMSCode))
	ug.POST("/login_sms", ginx.WrapReq[loginSMSReq](h.LoginSMS))
	ug.POST("/refresh_token", ginx.Wrap(h.RefreshToken))

	ug.GET("/profile", ginx.WrapClaim[jwtx.UserClaims](h.ProfileJWT))
	ug.POST("/edit", ginx.WrapReq[userEditReq](h.Edit))
}

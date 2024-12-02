package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/internal/service"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
	"net/http"
)

type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,len=10"`
}

func (h *UserHandler) LoginJWT(ctx *gin.Context, req loginReq) (any, error) {
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if errorx.IsCode(err, errcode.ErrInvalidUserNameOrPassword) {
		h.l.P3(ctx, "invalid username aor password", loggerx.Error(err))
		return nil, err
	}
	if err != nil {
		h.l.P1(ctx, "failed to login with jwt", loggerx.Error(err))
		return nil, err
	}

	if err = h.SetLoginToken(ctx, user.Id); err != nil {
		h.l.P1(ctx, "failed to set login token", loggerx.Error(err))
		return nil, err
	}

	return nil, nil
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 步骤2
	// 在这里登录成功了
	// 设置 session
	sess := sessions.Default(ctx)
	// 我可以随便设置值了
	// 你要放在 session 里面的值
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		// 一分钟过期
		MaxAge: 60,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

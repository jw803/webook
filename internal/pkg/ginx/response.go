package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/pkg/errorx"
	"net/http"
)

type Response struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		coder := errorx.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), Response{
			Code:    coder.Code(),
			Message: coder.String(),
			Data:    data,
		})
		return
	}
	c.JSON(http.StatusOK, Response{Data: data})
}

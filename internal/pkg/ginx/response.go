package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/pkg/errorx"
	"net/http"
)

type Response struct {
	Code      int         `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Reference string      `json:"reference,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		coder := errorx.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), Response{
			Code:    coder.Code(),
			Message: coder.Reference(),
			Data:    data,
		})
	}
	c.JSON(http.StatusOK, Response{Data: data})
}

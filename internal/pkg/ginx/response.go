package response

import (
	"github.com/gin-gonic/gin"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}

func SendResponse(ctx *gin.Context, info *Response, data any) {
	res := Result{}
	res.Code = info.Code()
	res.Msg = info.Msg()
	res.Data = data
	ctx.JSON(info.HttpStatusCode(), res)
}

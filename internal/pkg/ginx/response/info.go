package response

import (
	"fmt"
)

type Response struct {
	code int `json:"code"`

	httpCode int `json:"http_status_code"`

	msg string `json:"message"`
}

func NewResponse(code, httpCode int, msg string) *Response {
	return &Response{code: code, httpCode: httpCode, msg: msg}
}

func (e *Response) Error() string {
	return fmt.Sprintf("Error Code：%d, Error Message:：%s", e.Code(), e.Msg())
}

func (e *Response) Code() int {
	return e.code
}

func (e *Response) HttpStatusCode() int {
	return e.httpCode
}

func (e *Response) Msg() string {
	return e.msg
}

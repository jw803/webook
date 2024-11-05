package code

import "github.com/jw803/webook/pkg/errorx"

func init() {
	// common
	errorx.Register(ErrSystem, 500, errorx.System, "internal server error")
	errorx.Register(ErrBind, 400, errorx.System, "invalid json")
	errorx.Register(ErrValidation, 400, errorx.System, "failed to validation")
	errorx.Register(ErrTimeParsing, 400, errorx.Client, "failed to time string parsing")
	errorx.Register(ErrTokenInvalid, 401, errorx.System, "token invalid")

	// database errors
	errorx.Register(ErrDatabase, 500, errorx.System, "database error")

	// wstyle
	errorx.Register(ErrBadRequest, 400, errorx.System, "Bad Request")
	errorx.Register(ErrBadRequestFile, 400, errorx.System, "Bad File Format & Content")
}

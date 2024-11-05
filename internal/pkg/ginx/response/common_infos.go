package response

import (
	"context"
	"fmt"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/errorx"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/logging"
)

var l = logging.NewNoOpLogger()

func SetLogger(logger logging.Logger) {
	l = logger
}

func ResponseFactory(ctx context.Context, err error, alertLevel logging.AlertLevel, msgs ...string) *Response {
	coder := errorx.ParseCoder(err)
	if err != nil {
		logMsg := err.Error()
		if len(msgs) > 0 {
			logMsg = fmt.Sprintf("%s: %s", msgs[0], err.Error())
		}

		switch alertLevel {
		case logging.P0:
			l.Error(ctx, logging.P0, coder.ErrorType(), logMsg, logging.Error(err))
		case logging.P1:
			l.Error(ctx, logging.P1, coder.ErrorType(), logMsg, logging.Error(err))
		case logging.P2:
			l.Warn(ctx, logging.P2, coder.ErrorType(), logMsg, logging.Error(err))
		case logging.P3:
			l.Warn(ctx, logging.P3, coder.ErrorType(), logMsg, logging.Error(err))
		}
	}
	return NewResponse(coder.Code(), coder.HTTPStatus(), coder.External())
}

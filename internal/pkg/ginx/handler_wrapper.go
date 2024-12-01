package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/prometheus/client_golang/prometheus"
)

var logger = loggerx.NewNoOpLogger()

func SetLogger(l loggerx.Logger) {
	logger = l
}

var errorCounter *prometheus.CounterVec

func InitErrorCounter(opt prometheus.CounterOpts) {
	errorCounter = prometheus.NewCounterVec(opt, []string{"method", "path"})
	prometheus.MustRegister(errorCounter)
}

func Wrap(fn func(*gin.Context) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, fnErr := fn(ctx)
		WriteResponse(ctx, fnErr, data)
		return
	}
}

func WrapClaim[C jwt.Claims](fn func(*gin.Context, *C) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, ok := parseClaim[C](ctx)
		if !ok {
			return
		}

		data, fnErr := fn(ctx, claims)
		WriteResponse(ctx, fnErr, data)
		return
	}
}

func WrapReq[DTO any](fn func(ctx *gin.Context, dto DTO) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dto, ok := bindAndValidateReq[DTO](ctx)
		if ok {
			return
		}

		data, fnErr := fn(ctx, *dto)
		WriteResponse(ctx, fnErr, data)
		return
	}
}

func WrapQuery[Query any](fn func(*gin.Context, Query) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query, ok := bindAndValidateQuery[Query](ctx)
		if !ok {
			return
		}

		data, fnErr := fn(ctx, *query)
		WriteResponse(ctx, fnErr, data)
		return
	}
}

func WrapClaimsReq[C jwt.Claims, DTO any](fn func(*gin.Context, *C, DTO) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, ok := parseClaim[C](ctx)
		if !ok {
			return
		}

		dto, ok := bindAndValidateReq[DTO](ctx)
		if ok {
			return
		}

		data, fnErr := fn(ctx, claims, *dto)
		WriteResponse(ctx, fnErr, data)
		return
	}
}

func WrapClaimsQuery[C jwt.Claims, Query any](fn func(*gin.Context, *C, Query) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query, ok := bindAndValidateQuery[Query](ctx)
		if !ok {
			return
		}

		claims, ok := parseClaim[C](ctx)
		if !ok {
			return
		}

		data, fnErr := fn(ctx, claims, *query)
		WriteResponse(ctx, fnErr, data)
		return
	}
}

func WrapClaimsParam[C jwt.Claims, Param any](fn func(*gin.Context, *C, Param) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params, ok := bindAndValidateParams[Param](ctx)
		if !ok {
			return
		}

		claims, ok := parseClaim[C](ctx)
		if !ok {
			return
		}

		data, fnErr := fn(ctx, claims, *params)
		WriteResponse(ctx, fnErr, data)
		return
	}
}

func isValidationError(err error) bool {
	var invalidValidationError *validator.InvalidValidationError
	return errors.As(err, &invalidValidationError)
}

func parseClaim[C jwt.Claims](ctx *gin.Context) (*C, bool) {
	rawVal, ok := ctx.Get("claims")
	if !ok {
		logger.P3(ctx, "claim missing")
		WriteResponse(ctx, errorx.WithCode(errcode.ErrTokenMissing, "claim missing"), nil)
		return nil, false
	}

	claims, ok := rawVal.(*C)
	if !ok {
		logger.P3(ctx, "invalid token format")
		WriteResponse(ctx, errorx.WithCode(errcode.ErrTokenInvalid, "invalid token format"), nil)
		return nil, false
	}

	return claims, true
}

func bindAndValidateReq[DTO any](ctx *gin.Context) (*DTO, bool) {
	var dto DTO
	if err := ctx.ShouldBind(&dto); err != nil {
		if isValidationError(err) {
			logger.P3(ctx, "failed to validate request body", loggerx.Error(err))
			WriteResponse(ctx, errorx.WithCode(errcode.ErrValidation, "failed to validate request body"), nil)
			return nil, true
		}
		logger.P3(ctx, "failed to bind request body", loggerx.Error(err))
		WriteResponse(ctx, errorx.WithCode(errcode.ErrBind, "failed to bind request body"), nil)
		return nil, true
	}
	return &dto, false
}

func bindAndValidateQuery[DTO any](ctx *gin.Context) (*DTO, bool) {
	var dto DTO
	if err := ctx.ShouldBindQuery(&dto); err != nil {
		if isValidationError(err) {
			logger.P3(ctx, "failed to validate request query", loggerx.Error(err))
			WriteResponse(ctx, errorx.WithCode(errcode.ErrValidation, "failed to validate request query"), nil)
			return nil, true
		}
		logger.P3(ctx, "failed to bind request query", loggerx.Error(err))
		WriteResponse(ctx, errorx.WithCode(errcode.ErrBind, "failed to bind request query"), nil)
		return nil, true
	}
	return &dto, false
}

func bindAndValidateParams[DTO any](ctx *gin.Context) (*DTO, bool) {
	var dto DTO
	if err := ctx.ShouldBindUri(&dto); err != nil {
		if isValidationError(err) {
			logger.P3(ctx, "failed to validate request path params", loggerx.Error(err))
			WriteResponse(ctx, errorx.WithCode(errcode.ErrValidation, "failed to validate request path params"), nil)
			return nil, true
		}
		logger.P3(ctx, "failed to bind request body", loggerx.Error(err))
		WriteResponse(ctx, errorx.WithCode(errcode.ErrBind, "failed to bind request path params"), nil)
		return nil, true
	}
	return &dto, false
}

package ginx

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/jw803/webook/pkg/errorx"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/jw803/webook/pkg/ginx/response"
	"github.com/jw803/webook/pkg/logging"
)

var logger = logging.NewNoOpLogger()

func SetLogger(l logging.Logger) {
	logger = l
}

var errorCounter *prometheus.CounterVec

func InitErrorCounter(opt prometheus.CounterOpts) {
	errorCounter = prometheus.NewCounterVec(opt, []string{"method", "path"})
	prometheus.MustRegister(errorCounter)
}

func Wrap(fn func(*gin.Context) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, res := fn(ctx)
		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
		}
		response.SendResponse(ctx, res, data)
	}
}

func WrapClaim[T any](fn func(*gin.Context, *T) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawVal, ok := ctx.Get("claims")
		if !ok {
			logger.Error(ctx, "", "", "No Claim", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}
		claims, ok := rawVal.(*T)
		if !ok {
			logger.Error(ctx, "", "", "Token Format Error", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}

		data, res := fn(ctx, claims)

		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
			return
		}

		response.SendResponse(ctx, res, data)
	}
}

func WrapReq[DTO any](fn func(ctx *gin.Context, dto DTO) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto DTO
		if err := ctx.Bind(&dto); err != nil {
			logger.Error(ctx, "", "", "failed to parse request body", logging.Error(err))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}
		data, res := fn(ctx, dto)
		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
		}
		response.SendResponse(ctx, res, data)
	}
}

func WrapForm[DTO any](fn func(ctx *gin.Context, file multipart.File, fileHeader *multipart.FileHeader, dto DTO) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto DTO

		file, fileHeader, err := ctx.Request.FormFile("file")
		if err != nil {
			logger.Error(ctx, "", "", "failed to parse request body", logging.Error(err))
			res := response.NewResponse(400001, http.StatusBadRequest, "Bad File Format & Content")
			response.SendResponse(ctx, res, nil)
			return
		}
		defer file.Close()

		bodyString := ctx.Request.FormValue("body")
		if bodyString == "" {
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}
		if err := json.Unmarshal([]byte(bodyString), &dto); err != nil {
			logger.Error(ctx, "", "", "failed to parse request body", logging.Error(err))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}

		data, res := fn(ctx, file, fileHeader, dto)
		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
			return
		}

		response.SendResponse(ctx, res, data)
	}
}

func WrapQuery[Query any](fn func(*gin.Context, Query) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query Query
		if err := ctx.ShouldBindQuery(&query); err != nil {
			logger.Error(ctx, "", "", "Failed To Parse Request Body", logging.Error(err))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}

		data, res := fn(ctx, query)

		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
			return
		}

		response.SendResponse(ctx, res, data)
	}
}

func WrapClaimForm[DTO any](fn func(ctx *gin.Context, shoplineClaim ShoplineClaims, file multipart.File, fileHeader *multipart.FileHeader, dto DTO) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto DTO

		file, fileHeader, err := ctx.Request.FormFile("file")
		if err != nil {
			logger.Error(ctx, "", "", "failed to parse request body", logging.Error(err))
			res := response.NewResponse(400001, http.StatusBadRequest, "Bad File Format & Content")
			response.SendResponse(ctx, res, nil)
			return
		}
		defer file.Close()

		const maxFileSize = 1.5 * 1024 * 1024 // 1.5MB in bytes
		if fileHeader.Size > maxFileSize {
			logger.Error(ctx, "", "", "file size exceeds 10MB limit", logging.Error(err))
			res := response.NewResponse(400002, http.StatusBadRequest, "File size exceeds 10MB limit")
			response.SendResponse(ctx, res, nil)
			return
		}

		bodyString := ctx.Request.FormValue("body")
		if bodyString == "" {
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}

		if err := json.Unmarshal([]byte(bodyString), &dto); err != nil {
			logger.Error(ctx, "", "", "failed to parse request body", logging.Error(err))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}

		validate := validator.New()
		if err := validate.Struct(dto); err != nil {
			logger.Error(ctx, "", "", "Request body field validation failed", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}

		rawVal, ok := ctx.Get("shopline")
		if !ok {
			logger.Error(ctx, "", "", "No Shopline Claim", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}
		claims, ok := rawVal.(ShoplineClaims)
		if !ok {
			logger.Error(ctx, "", "", "Token Format Error", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}

		data, res := fn(ctx, claims, file, fileHeader, dto)
		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
			return
		}

		response.SendResponse(ctx, res, data)
	}
}

func WrapClaimsReq[DTO any](fn func(*gin.Context, ShoplineClaims, DTO) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto DTO
		if err := ctx.BindJSON(&dto); err != nil {
			logger.Error(ctx, "", "", "Failed To Parse Request Body", logging.Error(err))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}
		validate := validator.New()
		if err := validate.Struct(dto); err != nil {
			logger.Warn(ctx, "", "", "Request body field validation failed", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}

		rawVal, ok := ctx.Get("shopline")
		if !ok {
			logger.Error(ctx, "", "", "No Shopline Claim", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}
		claims, ok := rawVal.(ShoplineClaims)
		if !ok {
			logger.Error(ctx, "", "", "Token Format Error", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}

		data, res := fn(ctx, claims, dto)
		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
			return
		}

		response.SendResponse(ctx, res, data)
	}
}

func WrapClaimsQuery[Query any](fn func(*gin.Context, ShoplineClaims, Query) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query Query
		if err := ctx.ShouldBindQuery(&query); err != nil {
			logger.Error(ctx, "", "", "Failed To Parse Request Body", logging.Error(err))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}

		rawVal, ok := ctx.Get("shopline")
		if !ok {
			logger.Error(ctx, "", "", "No Shopline Claim", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}
		claims, ok := rawVal.(ShoplineClaims)
		if !ok {
			logger.Error(ctx, "", "", "Token Format Error", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}

		data, res := fn(ctx, claims, query)

		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
			return
		}

		response.SendResponse(ctx, res, data)
	}
}

func WrapClaimsReqWithoutClaims[DTO any](fn func(*gin.Context, DTO) (any, *response.Response)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto DTO
		if err := ctx.BindJSON(&dto); err != nil {
			logger.Error(ctx, "", "", "Failed To Parse Request Body", logging.Error(err))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}
		validate := validator.New()
		if err := validate.Struct(dto); err != nil {
			logger.Error(ctx, "", "", "Request body field validation failed", logging.String("path", ctx.Request.URL.Path))
			res := response.NewResponse(400000, http.StatusBadRequest, "Bad Request")
			response.SendResponse(ctx, res, nil)
			return
		}

		data, res := fn(ctx, dto)
		if res.Code() != errorx.CodeSuccess {
			response.SendResponse(ctx, res, nil)
			return
		}

		response.SendResponse(ctx, res, data)
	}
}

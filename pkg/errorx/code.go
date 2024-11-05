package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
)

const (
	CodeSuccess = 0
	CodeUnknown = 1
)

var (
	successCoder = &ErrCode{CodeSuccess, http.StatusOK, "ok", System}
	unknownCoder = &ErrCode{CodeUnknown, http.StatusInternalServerError, "internal server error", System}
)

var codes = map[int]*ErrCode{}
var codeMux = &sync.Mutex{}

type ErrCode struct {
	// C refers to the code of the ErrCode.
	code int

	// HTTP status that should be used for the associated error code.
	httpStatus int

	// External (user) facing error text.
	external string

	// Ref specify the reference document.
	errorType ErrorType
}

// Code returns the integer code of ErrCode.
func (coder ErrCode) Code() int {
	return coder.code
}

// String implements stringer. String returns the external error message,
// if any.
func (coder ErrCode) External() string {
	return coder.external
}

// Reference returns the reference document.
func (coder ErrCode) ErrorType() ErrorType {
	return coder.errorType
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder ErrCode) HTTPStatus() int {
	if coder.httpStatus == 0 {
		return http.StatusInternalServerError
	}
	return coder.httpStatus
}

// nolint: unparam
func Register(code int, httpStatus int, errorType ErrorType, external string) {
	found := slices.Contains([]int{200, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}

	coder := &ErrCode{
		code:       code,
		httpStatus: httpStatus,
		external:   external,
		errorType:  errorType,
	}

	mustRegister(coder)
}

func mustRegister(coder *ErrCode) {
	if coder.Code() == 0 || coder.Code() == 1 {
		panic("code '0' is reserved as ErrUnknown error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

func ParseCoder(err error) *ErrCode {
	if err == nil {
		return successCoder
	}

	var v *withCode
	if errors.As(err, &v) {
		if coder, ok := codes[v.code]; ok {
			ext := coder.external
			if v.external != "" {
				ext = v.external
			}
			return &ErrCode{
				code:       coder.code,
				httpStatus: coder.httpStatus,
				external:   ext,
				errorType:  coder.errorType,
			}
		}
	}

	return unknownCoder
}

func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return false
}

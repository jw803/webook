package code

const (
	// ErrSystem - 500: Internal server error.
	ErrSystem int = iota + 100001

	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind

	// ErrValidation - 400: Validation failed.
	ErrValidation

	ErrTimeZoneLoading

	// ErrTimeParsing - 400: Validation failed.
	ErrTimeParsing

	ErrNumberParsing

	ErrObjectIdParsing

	// ErrTokenInvalid - 401: Token invalid.
	ErrTokenInvalid

	ErrOpenFile

	ErrExcelRead

	ErrCronjobPartialFailed
)

// common: database errors.
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
)

const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301
	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed
)

// AWS
const (
	// ErrS3GetFile - Failedd to get file from AWS S3
	ErrS3GetFile int = iota + 100400
	ErrSQSSendMessage
)

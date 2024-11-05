package errorx

type ErrorType string

const (
	// Client 來自Client的非法請求
	Client     ErrorType = "Client"
	System     ErrorType = "System"
	DB         ErrorType = "DB"
	OpenAPI    ErrorType = "OpenAPI"
	AWS        ErrorType = "AWS"
	ThirdParty ErrorType = "ThirdParty"
)

package errcode

// init register error codes defines in this source code to `github.com/marmotedu/errors`
func init() {
	register(ErrSuccess, 200, "OK")
	register(ErrSystem, 500, "Internal server error")
	register(ErrUnknown, 500, "Internal server error")
	register(ErrBind, 400, "Invalid request format")
	register(ErrValidation, 400, "Validation failed")
	register(ErrCookieMissing, 400, "Cookie missing")
	register(ErrTokenMissing, 401, "Token missing")
	register(ErrTokenInvalid, 401, "Token invalid")
	register(ErrSessionInvalid, 401, "Session invalid")
	register(ErrPageNotFound, 404, "Page not found")
	register(ErrDatabase, 500, "Internal server error")
	register(ErrRedis, 500, "Internal server error")
	register(ErrEncrypt, 401, "Error occurred while encrypting the user password")
	register(ErrSignatureInvalid, 401, "Signature is invalid")
	register(ErrExpired, 401, "Token expired")
	register(ErrInvalidAuthHeader, 401, "Invalid authorization header")
	register(ErrMissingHeader, 401, "The `Authorization` header was empty")
	register(ErrPasswordIncorrect, 401, "Password was incorrect")
	register(ErrPermissionDenied, 403, "Permission denied")
	register(ErrMaliciousUser, 403, "Permission denied")
	register(ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
	register(ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
	register(ErrInvalidJSON, 500, "Data is not valid JSON")
	register(ErrEncodingJSON, 500, "JSON data could not be encoded")
	register(ErrDecodingJSON, 500, "JSON data could not be decoded")
	register(ErrInvalidYaml, 500, "Data is not valid Yaml")
	register(ErrEncodingYaml, 500, "Yaml data could not be encoded")
	register(ErrDecodingYaml, 500, "Yaml data could not be decoded")

	// UserModule
	register(ErrInvalidUserNameOrPassword, 500, "invalid username or password")
	register(ErrUserNotFound, 404, "user not found")
	register(ErrUserDuplicated, 400, "user duplicate")
	register(ErrInvalidPassword, 400, "Invalid password format")
	register(ErrPasswordNotMatch, 400, "The password and the confirmation password do not match")
	register(ErrDuplicateEmailSignUp, 400, "email has already been registered")

	register(ErrArticleNotFound, 404, "article not found")

	register(ErrSMSCodeSendTooFrequently, 401, "send sms code too frequency")
	register(ErrSMSCodeInvalid, 400, "sms verification code is incorrect")

	register(ErrWeChatStateMismatch, 400, "state mismatch")
	register(ErrWeChatVerificationCodeInvalid, 400, "wecaht verification code invalid")
}

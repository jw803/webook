package errcode

// User Module
const (
	ErrUserNotFound int = iota + 110001
	ErrInvalidUserNameOrPassword
	ErrUserDuplicated

	ErrInvalidPassword
	ErrPasswordNotMatch

	ErrSMSCodeSendTooFrequently
	ErrSMSCodeInvalid

	ErrArticleNotFound
	ErrMaliciousUser
)

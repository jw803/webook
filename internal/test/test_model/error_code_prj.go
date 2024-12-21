package test_model

// User Module
const (
	ErrUserNotFound int = iota + 110001
	ErrInvalidUserNameOrPassword
	ErrUserDuplicated
	ErrDuplicateEmailSignUp
	ErrInvalidPassword
	ErrPasswordNotMatch

	ErrSMSCodeSendTooFrequently
	ErrSMSCodeInvalid

	ErrArticleNotFound
	ErrMaliciousUser
)

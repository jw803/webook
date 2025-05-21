package errcode

// User Module
const (
	ErrUserNotFound int = iota + 110001
	ErrInvalidUserNameOrPassword
	ErrUserDuplicated
	ErrDuplicateEmailSignUp
	ErrUserCacheKeyNotFound

	ErrInvalidPassword
	ErrPasswordNotMatch

	ErrSMSCodeSendTooFrequently
	ErrSMSCodeInvalid

	ErrArticleNotFound
	ErrMaliciousUser

	ErrWeChatStateMismatch
	ErrWeChatVerificationCodeInvalid
)

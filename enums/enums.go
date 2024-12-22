package enums

type ErrorCode string

// Кастомный текст ошибок
const (
	ErrUnprocessableEntity ErrorCode = "Expression is not valid"
	ErrInternalServerError ErrorCode = "Internal server error"
)

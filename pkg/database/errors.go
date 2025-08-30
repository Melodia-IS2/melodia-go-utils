package database

type SqlErrorType int

const (
	SqlErrorTypeNotFound SqlErrorType = iota
	SqlErrorTypeConflict
	SqlErrorTypeFkError
	SqlErrorTypeUnknown
)

type SqlError struct {
	ErrorType SqlErrorType
	Entity    *string
	Message   string
}

func NewSqlError(errorType SqlErrorType, entity *string, message string) *SqlError {
	return &SqlError{ErrorType: errorType, Entity: entity, Message: message}
}

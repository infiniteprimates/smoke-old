package db

import (
	"fmt"
)

type (
	dbErrorCode int

	dbError struct {
		Code    dbErrorCode
		Message string
	}
)

const (
	EntityNotFound dbErrorCode = iota
	EntityExists
	Unknown
)

func NewDbError(code dbErrorCode, format string, args ...interface{}) error {
	return &dbError{code, fmt.Sprintf(format, args...)}
}

func (e *dbError) Error() string {
	return fmt.Sprintf("DB Error: Code = '%d', Message = '%s'", e.Code, e.Message)
}

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
	DbNotFound dbErrorCode = iota
	DbUnknown
)

func NewDbError(code dbErrorCode, message string) error {
	return &dbError{code, message}
}

func (e *dbError) Error() string {
	return fmt.Sprintf("DB Error: Code = '%d', Message = '%s'", e.Code, e.Message)
}

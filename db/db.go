package db

import (
	"fmt"
)

type (
	dbErrorReason string

	dbError struct {
		Reason  dbErrorReason
		Message string
	}
)

const (
	EntityNotFound = "EntityNotFound"
	EntityExists   = "EntityExists"
	Unknown        = "Unknown"
)

func NewDbError(code dbErrorReason, format string, args ...interface{}) error {
	return &dbError{code, fmt.Sprintf(format, args...)}
}

func (e *dbError) Error() string {
	return fmt.Sprintf("DB Error: Reason = '%s', Message = '%s'", e.Reason, e.Message)
}

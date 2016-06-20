package db

import (
	"fmt"
)

type (
	dbErrorReason string

	dbError struct {
		reason  dbErrorReason
		message string
	}
)

const (
	EntityNotFound = "EntityNotFound"
	EntityExists   = "EntityExists"
	Unknown        = "Unknown"
)

func newDbError(code dbErrorReason, format string, args ...interface{}) error {
	return &dbError{code, fmt.Sprintf(format, args...)}
}

func (e *dbError) Error() string {
	return fmt.Sprintf("DB Error: Reason = '%s', Message = '%s'", e.reason, e.message)
}

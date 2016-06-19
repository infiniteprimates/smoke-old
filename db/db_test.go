package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDbError(t *testing.T) {
	err := newDbError(EntityNotFound, "%s", "Catastrophe")

	if assert.NotNil(t, err, "Error was nil.") {
		if dbErr, ok := err.(*dbError); !ok {
			assert.True(t, ok, "Error was not a dbError.")
		} else {
			assert.Equal(t, dbErrorReason(EntityNotFound), dbErr.reason, "Reason did not match.")
			assert.Equal(t, "Catastrophe", dbErr.message, "Message did not match.")
			assert.Contains(t, dbErr.Error(), EntityNotFound, "Error does not contain expected reason.")
			assert.Contains(t, dbErr.Error(), "Catastrophe", "Error does not contain expected message.")
		}
	}

}

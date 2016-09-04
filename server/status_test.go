package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"strconv"
)

func TestStatusResource_newStatus_Success(t *testing.T) {
	err := newStatus(http.StatusNotFound)

	if assert.Error(t, err, "Expected error not returned") {
		assert.Equal(t, http.StatusNotFound, err.(*smokeStatus).Code, "Invalid status.")
		assert.Equal(t, statusMessages[http.StatusNotFound], err.(*smokeStatus).Message, "Invalid message.")
	}
}

func TestStatusResource_newStatus_SuccessUnknown(t *testing.T) {
	err := newStatus(http.StatusContinue)

	if assert.Error(t, err, "Expected error not returned") {
		assert.Equal(t, http.StatusContinue, err.(*smokeStatus).Code, "Invalid status.")
		assert.Equal(t, unknownMessage, err.(*smokeStatus).Message, "Invalid message.")
	}
}

func TestStatusResource_newStatusWithMessage_Success(t *testing.T) {
	err := newStatusWithMessage(http.StatusExpectationFailed, "fmt %s", "msg")

	if assert.Error(t, err, "Expected error not returned") {
		assert.Equal(t, http.StatusExpectationFailed, err.(*smokeStatus).Code, "Invalid status.")
		assert.Equal(t, "fmt msg", err.(*smokeStatus).Message, "Invalid message.")
	}
}

func TestStatusResource_smokeStatusError_Success(t *testing.T) {
	err := newStatus(http.StatusNotFound)

	if assert.Error(t, err, "Expected error not returned") {
		assert.Contains(t, err.Error(), strconv.Itoa(http.StatusNotFound), "Error string does not contain code.")
		assert.Contains(t, err.Error(), statusMessages[http.StatusNotFound], "Error string does not contain message.")
	}

}


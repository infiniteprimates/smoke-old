package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	"github.com/stretchr/testify/assert"
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

func TestStatusResource_smokeErrorHandler_smokeStatus(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", strings.NewReader("body"))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	err := &smokeStatus{
		999,
		"message",
	}

	handler := smokeErrorHandler()

	handler(err, c)

	status := unmarshallSmokeStatus(res.Body.String())
	assert.Equal(t, 999, res.Status(), "Invalid status.")
	assert.Equal(t, 999, status.Code, "Invalid smokeStatus code.")
	assert.Equal(t, "message", status.Message, "Invalid smokeStatus message.")
}

func TestStatusResource_smokeErrorHandler_smokeStatusUnauthorized(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", strings.NewReader("body"))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	err := &smokeStatus{
		http.StatusUnauthorized,
		"message",
	}

	handler := smokeErrorHandler()

	handler(err, c)

	status := unmarshallSmokeStatus(res.Body.String())
	assert.Equal(t, "Basic realm=\"Smoke\"", res.Header().Get(echo.HeaderWWWAuthenticate), "Invalid header.")
	assert.Equal(t, http.StatusUnauthorized, res.Status(), "Invalid status.")
	assert.Equal(t, http.StatusUnauthorized, status.Code, "Invalid smokeStatus code.")
	assert.Equal(t, "message", status.Message, "Invalid smokeStatus message.")
}

func TestStatusResource_smokeErrorHandler_httpError(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", strings.NewReader("body"))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	err := &echo.HTTPError{
		999,
		"message",
	}

	handler := smokeErrorHandler()

	handler(err, c)

	status := unmarshallSmokeStatus(res.Body.String())
	assert.Equal(t, 999, res.Status(), "Invalid status.")
	assert.Equal(t, 999, status.Code, "Invalid HTTPError code.")
	assert.Equal(t, "message", status.Message, "Invalid HTTPError message.")
}

func TestStatusResource_smokeErrorHandler_error(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", strings.NewReader("body"))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	err := errors.New("Plain old error.")

	handler := smokeErrorHandler()

	handler(err, c)

	status := unmarshallSmokeStatus(res.Body.String())
	assert.Equal(t, http.StatusInternalServerError, res.Status(), "Invalid status.")
	assert.Equal(t, http.StatusInternalServerError, status.Code, "Invalid HTTPError code.")
	assert.Equal(t, statusMessages[http.StatusInternalServerError], status.Message, "Invalid HTTPError message.")
}

func unmarshallSmokeStatus(s string) smokeStatus {
	var status smokeStatus
	if err := json.Unmarshal([]byte(s), &status); err != nil {
		panic(err)
	}
	return status
}

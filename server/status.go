package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type (
	smokeStatus struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

var statusMessages = map[int]string{
	http.StatusForbidden:           "Forbidden",
	http.StatusInternalServerError: "Internal server error",
	http.StatusNotFound:            "Not Found",
	http.StatusUnauthorized:        "Unauthorized",
}

func newStatus(code int) error {
	return newStatusWithMessage(code, statusMessages[code])
}

func newStatusWithMessage(code int, format string, args ...interface{}) error {
	return &smokeStatus{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func smokeErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	defaultHttpErrorHandler := e.DefaultHTTPErrorHandler
	return func(err error, c echo.Context) {
		if status, ok := err.(*smokeStatus); ok {
			if status.Code == http.StatusUnauthorized {
				c.Response().Header().Set(echo.HeaderWWWAuthenticate, "Basic realm=\"Smoke\"")
			}

			if !c.Response().Committed() {
				c.JSON(status.Code, status)
			}
		} else {
			defaultHttpErrorHandler(err, c)
		}
	}
}

func (s *smokeStatus) Error() string {
	return fmt.Sprintf("Code = '%d', Message = '%s'", s.Code, s.Message)
}

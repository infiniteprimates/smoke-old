package util

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	Status struct {
		Code int `json:"code"`
		Message string `json:"message"`
	}
)
var statusMessages = map[int]string{
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusInternalServerError: "Internal server error",
	http.StatusForbidden:           "Forbidden",
}

func AbortWithStatus(ctx echo.Context, code int) {
	AbortWithStatusAndMessage(ctx, code, statusMessages[code])
}

func AbortWithStatusAndMessage(ctx echo.Context, code int, msg string) {
	if code == http.StatusUnauthorized {
		ctx.Response().Header().Set(echo.HeaderWWWAuthenticate, "Basic realm=\"Smoke\"")
	}

	ctx.JSON(code, &Status{
		Code:    code,
		Message: msg,
	})
}

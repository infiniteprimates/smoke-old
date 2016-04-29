package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var statusMessages = map[int]string {
	http.StatusUnauthorized: "Unauthorized",
	http.StatusInternalServerError: "Internal server error",
	http.StatusForbidden: "Forbidden",
}

func AbortWithStatus(ctx *gin.Context, code int) {
	AbortWithStatusAndMessage(ctx, code, statusMessages[code])
}

func AbortWithStatusAndMessage(ctx *gin.Context, code int, msg string) {
	if code == http.StatusUnauthorized {
		ctx.Header("WWW-Authenticate", "Basic realm=\"Smoke\"")
	}

	ctx.JSON(code, gin.H{
		"code": code,
		"message": msg,
	})
	ctx.Abort()
}
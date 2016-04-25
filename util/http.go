package util

import "github.com/gin-gonic/gin"

var statusMessages = map[int]string {
	401: "Unauthorized",
	500: "Internal server error",
}

func AbortWithStatus(ctx *gin.Context, code int) {
	AbortWithStatusAndMessage(ctx, code, statusMessages[code])
}

func AbortWithStatusAndMessage(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"code": code,
		"message": msg,
	})
	ctx.Abort()
}
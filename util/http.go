package util

import "github.com/gin-gonic/gin"

func AbortWithStatus(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"code": code,
		"message": msg,
	})
	ctx.Abort()
}
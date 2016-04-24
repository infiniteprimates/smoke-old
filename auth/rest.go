package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/metrics"
	"github.com/infiniteprimates/smoke/util"
)

func CreateAuthResources(router gin.IRouter, middleware ...gin.HandlerFunc) {
	accounts := gin.Accounts{
		"admin": "secret",
		"user": "1234",
	}
	router.POST("/auth", metrics.MetricsHandler("get_auth"), gin.BasicAuth(accounts), postAuthorizationResource)
}

func postAuthorizationResource(ctx *gin.Context) {
	user := ctx.MustGet(gin.AuthUserKey).(string)

	if token, err := generateJwt(user) ; err != nil {
		util.AbortWithStatus(ctx, http.StatusInternalServerError, "Unknown error generating token")
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"type": "bearer",
			"token": token,
		})
	}
}

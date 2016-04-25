package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/metrics"
	"github.com/infiniteprimates/smoke/user"
	"github.com/infiniteprimates/smoke/util"
)

func CreateAuthResources(router gin.IRouter, middleware ...gin.HandlerFunc) {
	router.POST("/auth", metrics.MetricsHandler("get_auth"), postAuthorizationResource)
}

func postAuthorizationResource(ctx *gin.Context) {
	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		util.AbortWithStatus(ctx, http.StatusUnauthorized)
		return
	}

	user, err := user.Find(username)
	if err != nil {
		util.AbortWithStatus(ctx, http.StatusUnauthorized)
		return
	}

	if validatePassword(password, user.Password) {
		if token, err := generateJwt(user) ; err != nil {
			util.AbortWithStatusAndMessage(ctx, http.StatusInternalServerError, "Unknown error authorizing user")
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"type": "bearer",
				"token": token,
			})
			return
		}
	} else {
		util.AbortWithStatus(ctx, http.StatusUnauthorized)
		return
	}

}


package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/metrics"
)

func CreateUserResources(router gin.IRouter) {
	router.GET("/user", metrics.MetricsHandler("get_users"), getUsersResource)
	router.GET("/user/:userid", metrics.MetricsHandler("get_user"), getUserResource)
}

func getUsersResource(ctx *gin.Context) {
	ctx.String(http.StatusOK, "List of users")
}

func getUserResource(ctx *gin.Context) {
	userId := ctx.Param("userid")
	ctx.String(http.StatusOK, "User %s", userId)
}

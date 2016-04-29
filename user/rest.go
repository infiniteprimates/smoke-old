package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/auth"
	"github.com/infiniteprimates/smoke/metrics"
	"github.com/infiniteprimates/smoke/util"
)

func CreateUserResources(router gin.IRouter) {
	router.GET("/user", metrics.MetricsHandler("get_users"), auth.AuthorizationMiddleware(false), getUsersResource)
	router.GET("/user/:userid", metrics.MetricsHandler("get_user"), auth.AuthorizationMiddleware(false), getUserResource)
}

func getUsersResource(ctx *gin.Context) {
	if users, err := List() ; err != nil {
		util.AbortWithStatus(ctx, http.StatusInternalServerError)
	} else {
		ctx.JSON(http.StatusOK, users)
	}
}

func getUserResource(ctx *gin.Context) {
	userId := ctx.Param("userid")
	if user, err := Find(userId) ; err != nil {
		util.AbortWithStatus(ctx, http.StatusInternalServerError)
	} else if user == nil {
		util.AbortWithStatus(ctx, 404)
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}

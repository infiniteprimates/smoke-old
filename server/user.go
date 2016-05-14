package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/db"
	mw "github.com/infiniteprimates/smoke/middleware"
	"github.com/infiniteprimates/smoke/util"
)

func createUserResources(db *db.Db, router gin.IRouter) {
	router.GET("/user", mw.MetricsHandler("get_users"), mw.AuthorizationMiddleware(false), getUsersResource(db))
	router.GET("/user/:userid", mw.MetricsHandler("get_user"), mw.AuthorizationMiddleware(false), getUserResource(db))
}

func getUsersResource(db *db.Db) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if users, err := db.ListUsers(); err != nil {
			util.AbortWithStatus(ctx, http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, users)
		}
	}
}

func getUserResource(db *db.Db) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userid")
		if user, err := db.FindUser(userId); err != nil {
			util.AbortWithStatus(ctx, http.StatusInternalServerError)
		} else if user == nil {
			util.AbortWithStatus(ctx, 404)
		} else {
			ctx.JSON(http.StatusOK, user)
		}
	}
}

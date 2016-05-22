package server

import (
	"net/http"

	"github.com/infiniteprimates/smoke/db"
	"github.com/infiniteprimates/smoke/util"
	"github.com/labstack/echo"
)

func createUserResources(db *db.Db, group *echo.Group) {
	group.GET("/user", getUsersResource(db), metricsHandler("get_users"), authorizationMiddleware())
	group.GET("/user/:userid", getUserResource(db), metricsHandler("get_user"), authorizationMiddleware())
}

func getUsersResource(db *db.Db) echo.HandlerFunc {
	return func(c echo.Context) error {
		if users, err := db.ListUsers(); err != nil {
			return err
		} else {
			c.JSON(http.StatusOK, users)
		}

		return nil
	}
}

func getUserResource(db *db.Db) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		if user, err := db.FindUser(userId); err != nil {
			return err
		} else if user == nil {
			util.AbortWithStatus(c, 404)
		} else {
			c.JSON(http.StatusOK, user)
		}

		return nil
	}
}

package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/model"
	"github.com/infiniteprimates/smoke/service"
	"github.com/labstack/echo"
)

func createUserResources(group *echo.Group, cfg *config.Config, userService *service.UserService) {
	jwtKey := cfg.GetString(config.JwtKey)

	group.POST("/user", createUserResource(userService), metricsHandler("post_user"), authorizationMiddleware(jwtKey), requireAdminMiddleware("Only admins may create users."))
	group.GET("/user", getUsersResource(userService), metricsHandler("get_users"), authorizationMiddleware(jwtKey))
	group.GET("/user/:userid", getUserResource(userService), metricsHandler("get_user"), authorizationMiddleware(jwtKey))
	group.PUT("/user/:userid", updateUserResource(userService), metricsHandler("update_user"), authorizationMiddleware(jwtKey))
	group.DELETE("/user/:userid", deleteUserResource(userService), metricsHandler("delete_user"), authorizationMiddleware(jwtKey), requireAdminMiddleware("Only admins may delete users."))
	//group.PUT("/user/:userid/password", updateUserPasswordResource(userService), metricsHandler("update_user_password", authorizationMiddleware(jwtKey)))
}

func createUserResource(s *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(model.User)

		if err := c.Bind(user); err != nil {
			return newStatusWithMessage(http.StatusBadRequest, err.Error())
		}

		user, err := s.Create(user)
		if err != nil {
			return newStatusWithMessage(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusCreated, user)
	}
}

func getUserResource(s *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		user, err := s.Find(userId)
		if err != nil {
			return newStatus(404)
		}

		return c.JSON(http.StatusOK, user)
	}
}

func getUsersResource(s *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := s.List()
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, users)
	}
}

func updateUserResource(s *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		authUser := c.Get("user").(*jwt.Token).Claims["sub"].(string)
		isAdmin := c.Get("user").(*jwt.Token).Claims["isAdmin"].(bool)

		user := new(model.User)
		if err := c.Bind(user); err != nil {
			return newStatusWithMessage(http.StatusBadRequest, err.Error())
		}

		if user.Username != userId {
			return newStatusWithMessage(http.StatusBadRequest, "Url userId '%s' and json userId '%s' are mismatched.", userId, user.Username)
		}

		if authUser != userId && !isAdmin {
			return newStatusWithMessage(http.StatusForbidden, "Only admins may update other users.")
		}

		if !isAdmin && user.IsAdmin {
			return newStatusWithMessage(http.StatusForbidden, "Only admins may make other users admins.")
		}

		if isAdmin && !user.IsAdmin && authUser == user.Username {
			return newStatusWithMessage(http.StatusBadRequest, "For your own safety, I'm not going to allow you to remove admin privileges from yourself.")
		}

		user, err := s.Update(user)
		if err != nil {
			return newStatusWithMessage(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, user)
	}
}

func deleteUserResource(s *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		if authUser := c.Get("user").(*jwt.Token).Claims["sub"]; authUser == userId {
			return newStatusWithMessage(http.StatusForbidden, "You can't delete yourself.")
		}

		if err := s.Delete(userId); err != nil {
			return err
		}

		return newStatus(http.StatusNoContent)
	}
}

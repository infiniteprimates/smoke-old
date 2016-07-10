package server

import (
	"net/http"

	"github.com/infiniteprimates/smoke/model"
	"github.com/infiniteprimates/smoke/service"
	"github.com/labstack/echo"
)

func createUserResources(r router, authMiddleware echo.MiddlewareFunc, userService service.UserService) {
	group := r.Group("/users", authMiddleware)

	group.POST("", createUserResource(userService), metricsMiddleware("post_user"), requireAdminMiddleware("Only admins may create users."))
	group.GET("", getUsersResource(userService), metricsMiddleware("get_users"))
	group.GET("/:userid", getUserResource(userService), metricsMiddleware("get_user"))
	group.PUT("/:userid", updateUserResource(userService), metricsMiddleware("update_user"))
	group.DELETE("/:userid", deleteUserResource(userService), metricsMiddleware("delete_user"), requireAdminMiddleware("Only admins may delete users."))
	group.PUT("/:userid/password", updateUserPasswordResource(userService), metricsMiddleware("update_user_password"))
}

func createUserResource(s service.UserService) echo.HandlerFunc {
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

func getUserResource(s service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		user, err := s.Find(userId)
		if err != nil {
			return newStatus(404)
		}

		return c.JSON(http.StatusOK, user)
	}
}

func getUsersResource(s service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := s.List()
		if err != nil {
			return newStatusWithMessage(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, users)
	}
}

func updateUserResource(s service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		claims := extractClaims(c)
		authUser := claims["sub"].(string)
		isAdmin := claims["isAdmin"].(bool)

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

func deleteUserResource(s service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		claims := extractClaims(c)
		if authUser := claims["sub"]; authUser == userId {
			return newStatusWithMessage(http.StatusForbidden, "You can't delete yourself.")
		}

		if err := s.Delete(userId); err != nil {
			return err
		}

		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}
}

func updateUserPasswordResource(userService service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		claims := extractClaims(c)
		authUser := claims["sub"].(string)
		isAdmin := claims["isAdmin"].(bool)

		passwordReset := new(model.PasswordReset)
		if err := c.Bind(passwordReset); err != nil {
			return newStatusWithMessage(http.StatusBadRequest, err.Error())
		}

		if authUser != userId && !isAdmin {
			return newStatusWithMessage(http.StatusForbidden, "Only admins may set other user's passwords.")
		}

		if err := userService.UpdateUserPassword(userId, passwordReset, isAdmin); err != nil {
			return newStatusWithMessage(http.StatusBadRequest, "Password reset failed.")
		}

		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}
}

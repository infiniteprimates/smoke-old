package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/infiniteprimates/smoke/model"
	"github.com/infiniteprimates/smoke/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func createAuthResources(r router, authService service.AuthService) {
	r.POST("/auth", postAuthorizationResource(authService), metricsMiddleware("get_auth"), basicAuthExtractor())
}

func postAuthorizationResource(authService service.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		password := c.Get("password").(string)

		token, err := authService.AuthenticateUser(username, password)
		if err != nil {
			return newStatus(http.StatusUnauthorized)
		}

		return c.JSON(http.StatusOK, &model.Auth{
			AuthType: "bearer",
			Token:    token,
		})
	}
}

func authorizationMiddleware(jwtKey string) echo.MiddlewareFunc {
	return middleware.JWT([]byte(jwtKey))
}

func requireAdminMiddleware(message string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			if !user.Claims["isAdmin"].(bool) {
				return newStatusWithMessage(http.StatusForbidden, message)
			}
			return next(c)
		}
	}
}

func basicAuthExtractor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return middleware.BasicAuth(func(username string, password string) bool {
				c.Set("username", username)
				c.Set("password", password)
				return true
			})(next)(c)
		}
	}
}

package server

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/infiniteprimates/smoke/model"
	"github.com/infiniteprimates/smoke/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const ( //TODO: Put key into config
	Issuer       = "Smoke"
	JwtKey       = "s3kr1t"
	MethodBearer = "Bearer"
)

type (
	authResponse struct {
		AuthType string `json:"type"`
		Token    string `json:"token"`
	}
)

func createAuthResources(userService *service.UserService, passwordService *service.PasswordService, group *echo.Group) {
	group.POST("/auth", postAuthorizationResource(userService, passwordService), metricsHandler("get_auth"), basicAuthExtractor())
}

func postAuthorizationResource(userService *service.UserService, passwordService *service.PasswordService) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		password := c.Get("password").(string)

		user, err := userService.Find(username, true)
		if err != nil {
			return newStatus(http.StatusUnauthorized)
		}

		if !passwordService.ValidatePassword(password, user.Password) {
			return newStatus(http.StatusUnauthorized)
		}

		token, err := generateJwt(user)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &authResponse{
			AuthType: "bearer",
			Token:    token,
		})
	}
}

func generateJwt(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iss"] = Issuer
	token.Claims["sub"] = user.Username
	token.Claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	token.Claims["isAdmin"] = user.IsAdmin

	return token.SignedString([]byte(JwtKey))
}

func authorizationMiddleware() echo.MiddlewareFunc {
	return middleware.JWT([]byte(JwtKey))
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

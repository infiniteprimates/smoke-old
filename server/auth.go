package server

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/infiniteprimates/smoke/db"
	"github.com/infiniteprimates/smoke/model"
	"github.com/infiniteprimates/smoke/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
)

const ( //TODO: Put key into config
	ISSUER  = "Smoke"
	JWT_KEY = "s3kr1t"
	METHOD_BEARER = "Bearer"
)

type (
	authResponse struct {
		AuthType string `json:"type"`
		Token string `json:"token"`
	}
)

func createAuthResources(db *db.Db, group *echo.Group) {
	group.POST("/auth", postAuthorizationResource(db), metricsHandler("get_auth"), basicAuthExtractor())
}

func postAuthorizationResource(db *db.Db) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		password := c.Get("password").(string)

		user, err := db.FindUser(username)
		if err != nil {
			util.AbortWithStatus(c, http.StatusUnauthorized)
			return nil
		}

		if !validatePassword(password, user.Password) {
			util.AbortWithStatus(c, http.StatusUnauthorized)
			return nil
		}

		token, err := generateJwt(user)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, &authResponse{
			AuthType:  "bearer",
			Token: token,
		})

		return nil
	}
}

func validatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJwt(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iss"] = ISSUER
	token.Claims["sub"] = user.Username
	token.Claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	token.Claims["isAdmin"] = user.IsAdmin

	return token.SignedString([]byte(JWT_KEY))
}

func authorizationMiddleware() echo.MiddlewareFunc {
	return middleware.JWT([]byte(JWT_KEY))
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
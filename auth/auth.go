package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/user"
	"github.com/infiniteprimates/smoke/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	ISSUER = "Smoke"
	JWT_KEY = "s3kr1t"
)
func AuthorizationMiddleware(requireAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := jwt.ParseFromRequest(ctx.Request, keyFunc)

		if err != nil || !token.Valid {
			util.AbortWithStatus(ctx, http.StatusForbidden)
			return
		}

		ctx.Set("username", token.Claims["sub"])
		ctx.Set("isAdmin", token.Claims["isAdmin"])

		ctx.Next()
	}
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHS256) ; ok {
		return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
	}

	return []byte("s3kr1t"), nil
}

func validatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJwt(user *user.User) (string,error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iss"] = ISSUER
	token.Claims["sub"] = user.Username
	token.Claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	token.Claims["isAdmin"] = user.IsAdmin

	return token.SignedString([]byte(JWT_KEY))
}

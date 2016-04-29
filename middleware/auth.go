package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/util"
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
	if method, ok := token.Method.(*jwt.SigningMethodHMAC) ; !ok || method.Name != "HS256" {
		return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
	}

	return []byte("s3kr1t"), nil
}

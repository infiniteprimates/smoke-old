package rest

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/db"
	mw "github.com/infiniteprimates/smoke/middleware"
	"github.com/infiniteprimates/smoke/model"
	"github.com/infiniteprimates/smoke/util"
	"golang.org/x/crypto/bcrypt"
)

const (
	ISSUER = "Smoke"
	JWT_KEY = "s3kr1t"
)

func CreateAuthResources(router gin.IRouter) {
	router.POST("/auth", mw.MetricsHandler("get_auth"), postAuthorizationResource)
}

func postAuthorizationResource(ctx *gin.Context) {
	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		util.AbortWithStatus(ctx, http.StatusUnauthorized)
		return
	}

	user, err := db.FindUser(username)
	if err != nil {
		util.AbortWithStatus(ctx, http.StatusUnauthorized)
		return
	}

	if !validatePassword(password, user.Password) {
		util.AbortWithStatus(ctx, http.StatusUnauthorized)
		return
	}

	token, err := generateJwt(user)
	if err != nil {
		util.AbortWithStatusAndMessage(ctx, http.StatusInternalServerError, "Unknown error authorizing user")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"type": "bearer",
		"token": token,
	})
}

func validatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJwt(user *model.User) (string,error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iss"] = ISSUER
	token.Claims["sub"] = user.Username
	token.Claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	token.Claims["isAdmin"] = user.IsAdmin

	return token.SignedString([]byte(JWT_KEY))
}

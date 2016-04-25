package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/infiniteprimates/smoke/user"
	"golang.org/x/crypto/bcrypt"
)

const (
	ISSUER = "Smoke"
)

func validatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJwt(user *user.User) (string,error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iss"] = ISSUER
	token.Claims["sub"] = user.Username
	token.Claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	token.Claims["admin"] = user.IsAdmin

	return token.SignedString([]byte("s3kr1t"))
}

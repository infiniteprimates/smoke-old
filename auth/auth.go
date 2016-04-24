package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	ISSUER = "Smoke"
)
func generateJwt(user string) (string,error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iss"] = ISSUER
	token.Claims["sub"] = user
	token.Claims["exp"] = time.Now().Add(1 * time.Hour).Unix()

	admin := false
	if (user == "admin") {
		admin = true
	}
	token.Claims["admin"] = admin

	return token.SignedString([]byte("s3kr1t"))
}

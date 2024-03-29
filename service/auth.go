package service

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/db"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthService interface {
		AuthenticateUser(username string, password string) (string, error)
		HashPassword(password string) (string, error)
	}

	authService struct {
		cfg    config.Config
		userDb db.UserDb
	}
)

const (
	Issuer = "Smoke"
)

func NewAuthService(cfg config.Config, userDb db.UserDb) AuthService {
	return &authService{
		cfg:    cfg,
		userDb: userDb,
	}
}

func (s *authService) AuthenticateUser(username string, password string) (string, error) {
	user, err := s.userDb.Find(username)
	if err != nil {
		// hash the password so this takes time like a validation
		_, _ = s.HashPassword(password)
		return "", errors.New("Invalid credentials.")
	}

	if !s.validatePassword(password, user.PasswordHash) {
		return "", errors.New("Invalid credentials.")
	}

	return s.generateJwt(user)
}

func (s *authService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (s *authService) validatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *authService) generateJwt(user *db.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["iss"] = Issuer
	claims["sub"] = user.Username
	//TODO: Configurable expiration
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims["isAdmin"] = user.IsAdmin

	return token.SignedString([]byte(s.cfg.GetString(config.JwtKey)))
}

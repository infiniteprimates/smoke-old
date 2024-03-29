package service

import (
	"errors"
	"testing"

	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/db"
	mockcfg "github.com/infiniteprimates/smoke/mocks/config"
	mockdb "github.com/infiniteprimates/smoke/mocks/db"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_AuthenticateUser_Success(t *testing.T) {
	cfg := new(mockcfg.ConfigMock)
	userDb := new(mockdb.UserDbMock)
	svc := NewAuthService(cfg, userDb)
	user := &db.User{
		Username: "user",
	}
	passwordHash, err := svc.HashPassword("password")
	if !assert.NoError(t, err, "Password hashing failed.") {
		return
	}

	user.PasswordHash = passwordHash

	userDb.On("Find", "username").Return(user, nil)
	cfg.On("GetString", config.JwtKey).Return(config.JwtKey)

	token, err := svc.AuthenticateUser("username", "password")

	userDb.AssertExpectations(t)
	cfg.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred authenticating user.") {
		assert.NotEmpty(t, token, "An empty token was returned.")
	}
}

func TestAuthService_AuthenticateUser_BadPassword(t *testing.T) {
	cfg := new(mockcfg.ConfigMock)
	userDb := new(mockdb.UserDbMock)
	svc := NewAuthService(cfg, userDb)
	user := &db.User{
		Username: "user",
	}
	passwordHash, err := svc.HashPassword("password")
	if !assert.NoError(t, err, "Password hashing failed.") {
		return
	}

	user.PasswordHash = passwordHash

	userDb.On("Find", "username").Return(user, nil)

	token, err := svc.AuthenticateUser("username", "badpassword")

	userDb.AssertExpectations(t)
	cfg.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Empty(t, token, "Token was not empty.")
	}
}

func TestAuthService_AuthenticateUser_NoUser(t *testing.T) {
	cfg := new(mockcfg.ConfigMock)
	userDb := new(mockdb.UserDbMock)
	svc := NewAuthService(cfg, userDb)

	userDb.On("Find", "username").Return(nil, errors.New("Failure"))

	token, err := svc.AuthenticateUser("username", "password")

	userDb.AssertExpectations(t)
	cfg.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Empty(t, token, "Token was not empty.")
	}
}

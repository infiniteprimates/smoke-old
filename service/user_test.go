package service

import (
	"errors"
	"testing"

	"github.com/infiniteprimates/smoke/db"
	mockdb "github.com/infiniteprimates/smoke/mocks/db"
	mockservice "github.com/infiniteprimates/smoke/mocks/service"
	"github.com/infiniteprimates/smoke/model"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create_Success(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	authSvc := new(mockservice.AuthServiceMock)
	svc := NewUserService(userDb, authSvc)

	user := model.User{
		Username: "username",
		IsAdmin: false,
	}
	userEntity := db.User{
		Username: user.Username,
		IsAdmin: user.IsAdmin,
	}

	userDb.On("Create", &userEntity).Return(&userEntity, nil)

	result, err := svc.Create(&user)
	if assert.NoError(t, err, "An error occured creating a user.") {
		assert.Equal(t, user.Username, result.Username, "Username does not match.")
		assert.Equal(t, user.IsAdmin, result.IsAdmin, "IsAdmin does not match.")
	}
}

func TestUserService_Create_Error(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	authSvc := new(mockservice.AuthServiceMock)
	svc := NewUserService(userDb, authSvc)

	user := model.User{
		Username: "username",
		IsAdmin: false,
	}
	userEntity := db.User{
		Username: user.Username,
		IsAdmin: user.IsAdmin,
	}

	userDb.On("Create", &userEntity).Return(nil, errors.New("Danger!!"))

	result, err := svc.Create(&user)
	if assert.Error(t, err, "Expected error not returned.") {
		assert.Nil(t, result, "User was not nil.")
	}
}

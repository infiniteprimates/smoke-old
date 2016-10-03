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
	svc := NewUserService(userDb, nil)

	user := model.User{
		Username: "username",
		IsAdmin:  false,
	}

	userEntity := db.User{
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	userDb.On("Create", &userEntity).Return(&userEntity, nil)

	result, err := svc.Create(&user)

	userDb.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred creating a user.") {
		assert.Equal(t, user.Username, result.Username, "Username does not match.")
		assert.Equal(t, user.IsAdmin, result.IsAdmin, "IsAdmin does not match.")
	}
}

func TestUserService_Create_Error(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	user := model.User{
		Username: "username",
		IsAdmin:  false,
	}

	userEntity := db.User{
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	userDb.On("Create", &userEntity).Return(nil, errors.New("Danger!!"))

	result, err := svc.Create(&user)

	userDb.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Nil(t, result, "User was not nil.")
	}
}

func TestUserService_Find_Success(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	userEntity := db.User{
		Username: "username",
		IsAdmin:  false,
	}

	userDb.On("Find", "username").Return(&userEntity, nil)

	result, err := svc.Find("username")

	userDb.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred creating a user.") {
		assert.Equal(t, userEntity.Username, result.Username, "Username does not match.")
		assert.Equal(t, userEntity.IsAdmin, result.IsAdmin, "IsAdmin does not match.")
	}
}

func TestUserService_Find_Error(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	userDb.On("Find", "username").Return(nil, errors.New("Danger!!"))

	result, err := svc.Find("username")

	userDb.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Nil(t, result, "User was not nil.")
	}
}

func TestUserService_List_Success(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	user := &model.User{
		Username: "username",
		IsAdmin:  false,
	}

	userEntities := []*db.User{{
		Username: "username",
		IsAdmin:  false,
	}}

	userDb.On("List").Return(userEntities, nil)

	result, err := svc.List()

	userDb.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred listing users.") {
		assert.Len(t, result, len(userEntities), "Result length is incorrect.")
		assert.Contains(t, result, user, "Expected user is not contained in the result.")
	}
}

func TestUserService_List_Error(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	userDb.On("List").Return(nil, errors.New("Failure"))

	result, err := svc.List()

	userDb.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Nil(t, result, "User list was not nil.")
	}
}

func TestUserService_Update_Success(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	user := model.User{
		Username: "username",
		IsAdmin:  false,
	}

	userEntity := db.User{
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	userDb.On("Update", &userEntity).Return(&userEntity, nil)

	result, err := svc.Update(&user)

	userDb.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred updating a user.") {
		assert.Equal(t, user.Username, result.Username, "Username does not match.")
		assert.Equal(t, user.IsAdmin, result.IsAdmin, "IsAdmin does not match.")
	}
}

func TestUserService_Update_Error(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	user := model.User{
		Username: "username",
		IsAdmin:  false,
	}

	userEntity := db.User{
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	userDb.On("Update", &userEntity).Return(nil, errors.New("Oh noes!"))

	result, err := svc.Update(&user)

	userDb.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Nil(t, result, "Result was not nil.")
	}
}

func TestUserService_Delete_Success(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	userDb.On("Delete", "username").Return(nil)

	err := svc.Delete("username")

	userDb.AssertExpectations(t)

	assert.NoError(t, err, "An error occurred creating a user.")
}

func TestUserService_Delete_Error(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	svc := NewUserService(userDb, nil)

	userDb.On("Delete", "username").Return(errors.New("Bad things."))

	err := svc.Delete("username")
	assert.Error(t, err, "Expected error not returned.")
}

func TestUserService_UpdateUserPassword_SuccessNonAdmin(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	authSvc := new(mockservice.AuthServiceMock)
	svc := NewUserService(userDb, authSvc)

	passwordReset := &model.PasswordReset{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	authSvc.On("AuthenticateUser", "username", "oldpassword").Return("token", nil)
	authSvc.On("HashPassword", "newpassword").Return("hashedpassword", nil)
	userDb.On("UpdateUserPassword", "username", "hashedpassword").Return(nil)

	err := svc.UpdateUserPassword("username", passwordReset, false)

	userDb.AssertExpectations(t)
	authSvc.AssertExpectations(t)

	assert.NoError(t, err, "An error occurred updating user password.")
}

func TestUserService_UpdateUserPassword_SuccessAdmin(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	authSvc := new(mockservice.AuthServiceMock)
	svc := NewUserService(userDb, authSvc)

	passwordReset := &model.PasswordReset{
		NewPassword: "newpassword",
	}

	authSvc.On("HashPassword", "newpassword").Return("hashedpassword", nil)
	userDb.On("UpdateUserPassword", "username", "hashedpassword").Return(nil)

	err := svc.UpdateUserPassword("username", passwordReset, true)

	userDb.AssertExpectations(t)
	authSvc.AssertExpectations(t)

	assert.NoError(t, err, "An error occurred updating user password.")
}

func TestUserService_UpdateUserPassword_ErrorNonAdminBadAuth(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	authSvc := new(mockservice.AuthServiceMock)
	svc := NewUserService(userDb, authSvc)

	passwordReset := &model.PasswordReset{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	authSvc.On("AuthenticateUser", "username", "oldpassword").Return("", errors.New("Bad auth"))

	err := svc.UpdateUserPassword("username", passwordReset, false)

	userDb.AssertExpectations(t)
	authSvc.AssertExpectations(t)

	assert.Error(t, err, "Expected error not returned.")
}

func TestUserService_UpdateUserPassword_ErrorWhitespacePassword(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	authSvc := new(mockservice.AuthServiceMock)
	svc := NewUserService(userDb, authSvc)

	passwordReset := &model.PasswordReset{
		NewPassword: " \t",
	}

	err := svc.UpdateUserPassword("username", passwordReset, true)

	userDb.AssertExpectations(t)
	authSvc.AssertExpectations(t)

	assert.Error(t, err, "Expected error not returned.")
}

func TestUserService_UpdateUserPassword_ErrorHashFailure(t *testing.T) {
	userDb := new(mockdb.UserDbMock)
	authSvc := new(mockservice.AuthServiceMock)
	svc := NewUserService(userDb, authSvc)

	passwordReset := &model.PasswordReset{
		NewPassword: "newpassword",
	}

	authSvc.On("HashPassword", "newpassword").Return("", errors.New("Hash failure."))

	err := svc.UpdateUserPassword("username", passwordReset, true)

	userDb.AssertExpectations(t)
	authSvc.AssertExpectations(t)

	assert.Error(t, err, "Expected error not returned.")
}

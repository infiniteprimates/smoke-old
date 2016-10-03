package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	mockserver "github.com/infiniteprimates/smoke/mocks/server"
	mockservice "github.com/infiniteprimates/smoke/mocks/service"
	"github.com/infiniteprimates/smoke/model"
	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

func TestUserResource_createUserResources(t *testing.T) {
	router := new(mockserver.RouterMock)
	userSvc := new(mockservice.UserServiceMock)
	// echo folks need to use more interfaces...
	e := echo.New()
	group := e.Group("/")

	router.On("Group", "/users", mock.AnythingOfType("[]echo.MiddlewareFunc")).Return(group)
	createUserResources(router, func(echo.HandlerFunc) echo.HandlerFunc { return nil }, userSvc)
}

func TestUserResource_createUserResource_Success(t *testing.T) {
	userSvc := new(mockservice.UserServiceMock)
	user := &model.User{
		Username: "fred",
		IsAdmin: false,
	}
	body := marshallModel(user)
	e := echo.New()
	req := test.NewRequest(echo.POST, "/users", strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	userSvc.On("Create", user).Return(user, nil)

	handler := createUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred invoking handler.") {
		assert.Equal(t, http.StatusCreated, res.Status(), "Invalid status.")
		assert.Equal(t, body, res.Body.String(), "Invalid response.")
	}
}

func TestUserResource_createUserResource_BadJson(t *testing.T) {
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.POST, "/users", strings.NewReader("bad json"))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	handler := createUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, err.(*smokeStatus).Code, http.StatusBadRequest, "Invalid status.")
	}
}

func TestUserResource_createUserResource_CreateError(t *testing.T) {
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	user := &model.User{
		Username: "fred",
		IsAdmin: false,
	}
	req := test.NewRequest(echo.POST, "/users", strings.NewReader(marshallModel(user)))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	userSvc.On("Create", user).Return(nil, errors.New("Failure"))

	handler := createUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, err.(*smokeStatus).Code, http.StatusInternalServerError, "Invalid status.")
	}
}

func TestUserResource_getUserResource_Success(t *testing.T) {
	userId := "barney"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.GET, "/users", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	c.SetParamNames("userid")
	c.SetParamValues(userId)
	user := &model.User{
		Username: userId,
		IsAdmin: false,
	}

	userSvc.On("Find", userId).Return(user, nil)

	handler := getUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.NoError(t, err, "Unexpected error.") {
		assert.Equal(t, http.StatusOK, res.Status(), "Invalid status.")
		assert.Equal(t, marshallModel(user), res.Body.String(), "Invalid body.")
	}
}

func TestUserResource_getUserResource_NotFound(t *testing.T) {
	userId := "barney"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.GET, "/users", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	c.SetParamNames("userid")
	c.SetParamValues(userId)

	userSvc.On("Find", userId).Return(nil, nil)

	handler := getUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusNotFound, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_getUserResource_Failure(t *testing.T) {
	userId := "barney"
	failureMsg := "Failure"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.GET, "/users", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	c.SetParamNames("userid")
	c.SetParamValues(userId)

	userSvc.On("Find", userId).Return(nil, errors.New(failureMsg))

	handler := getUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusInternalServerError, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_getUsersResource_Success(t *testing.T) {
	userSvc := new(mockservice.UserServiceMock)
	users := []*model.User{
		{
			Username: "fred",
			IsAdmin: false,
		},
	}
	e := echo.New()
	req := test.NewRequest(echo.GET, "/users", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	userSvc.On("List").Return(users, nil)

	handler := getUsersResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred invoking handler.") {
		assert.Equal(t, http.StatusOK, res.Status(), "Invalid status.")
		assert.Equal(t, marshallModel(users), res.Body.String(), "Invalid response.")
	}
}

func TestUserResource_getUsersResource_Failure(t *testing.T) {
	failureMsg := "Failure"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.GET, "/users", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	userSvc.On("List").Return(nil, errors.New(failureMsg))

	handler := getUsersResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusInternalServerError, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_updateUserResource_SuccessUser(t *testing.T) {
	username := "bambam"
	userSvc := new(mockservice.UserServiceMock)
	user := &model.User{
		Username: username,
		IsAdmin: false,
	}
	body := marshallModel(user)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("Update", user).Return(user, nil)

	handler := updateUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred invoking handler.") {
		assert.Equal(t, http.StatusOK, res.Status(), "Invalid status.")
		assert.Equal(t, body, res.Body.String(), "Invalid response.")
	}
}

func TestUserResource_updateUserResource_SuccessAdmin(t *testing.T) {
	username := "bambam"
	userSvc := new(mockservice.UserServiceMock)
	user := &model.User{
		Username: username,
		IsAdmin: false,
	}
	body := marshallModel(user)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "someAdmin"
	claims["isAdmin"] = true
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("Update", user).Return(user, nil)

	handler := updateUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred invoking handler.") {
		assert.Equal(t, http.StatusOK, res.Status(), "Invalid status.")
		assert.Equal(t, body, res.Body.String(), "Invalid response.")
	}
}

func TestUserResource_updateUserResource_UpdateFailed(t *testing.T) {
	username := "bambam"
	failureMsg := "Failed"
	userSvc := new(mockservice.UserServiceMock)
	user := &model.User{
		Username: username,
		IsAdmin: false,
	}
	body := marshallModel(user)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("Update", user).Return(nil, errors.New(failureMsg))

	handler := updateUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusInternalServerError, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_updateUserResource_BadJson(t *testing.T) {
	username := "bambam"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader("Bad JSON"))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	handler := updateUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusBadRequest, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_updateUserResource_UrlAndModelMismatch(t *testing.T) {
	username := "bambam"
	userSvc := new(mockservice.UserServiceMock)
	user := &model.User{
		Username: username,
		IsAdmin: false,
	}
	body := marshallModel(user)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/bogus", strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues("bogus")

	handler := updateUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusBadRequest, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_updateUserResource_NotAdminAndNotSelf(t *testing.T) {
	username := "bambam"
	userSvc := new(mockservice.UserServiceMock)
	user := &model.User{
		Username: username,
		IsAdmin: false,
	}
	body := marshallModel(user)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "not" + username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	handler := updateUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusForbidden, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_updateUserResource_NotAdminMakeSelfAdmin(t *testing.T) {
	username := "bambam"
	userSvc := new(mockservice.UserServiceMock)
	user := &model.User{
		Username: username,
		IsAdmin: true,
	}
	body := marshallModel(user)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	handler := updateUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusForbidden, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_updateUserResource_DemoteSelf(t *testing.T) {
	username := "bambam"
	userSvc := new(mockservice.UserServiceMock)
	user := &model.User{
		Username: username,
		IsAdmin: false,
	}
	body := marshallModel(user)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = true
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	handler := updateUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusBadRequest, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_deleteUserResource_Success(t *testing.T) {
	username := "pebbles"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.DELETE, "/users/" + username, strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "not" + username
	claims["isAdmin"] = true
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("Delete", username).Return(nil)

	handler := deleteUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred invoking handler.") {
		assert.Equal(t, http.StatusNoContent, res.Status(), "Invalid status.")
	}
}

func TestUserResource_deleteUserResource_Failure(t *testing.T) {
	username := "pebbles"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.DELETE, "/users/" + username, strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "not" + username
	claims["isAdmin"] = true
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("Delete", username).Return(errors.New("Failure"))

	handler := deleteUserResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusInternalServerError, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_deleteUserResource_DeleteSelf(t *testing.T) {
	username := "pebbles"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.DELETE, "/users/" + username, strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = true
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("Delete", username).Return(nil)

	handler := deleteUserResource(userSvc)

	err := handler(c)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusForbidden, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_updateUserPasswordResource_SuccessUser(t *testing.T) {
	username := "betty"
	userSvc := new(mockservice.UserServiceMock)
	passwordReset := &model.PasswordReset{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}
	body := marshallModel(passwordReset)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("UpdateUserPassword", username, passwordReset, false).Return(nil)

	handler := updateUserPasswordResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred invoking handler.") {
		assert.Equal(t, http.StatusNoContent, res.Status(), "Invalid status.")
	}
}

func TestUserResource_updateUserPasswordResource_SuccessAdmin(t *testing.T) {
	username := "betty"
	userSvc := new(mockservice.UserServiceMock)
	passwordReset := &model.PasswordReset{
		NewPassword: "newpassword",
	}
	body := marshallModel(passwordReset)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = true
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("UpdateUserPassword", username, passwordReset, true).Return(nil)

	handler := updateUserPasswordResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.NoError(t, err, "An error occurred invoking handler.") {
		assert.Equal(t, http.StatusNoContent, res.Status(), "Invalid status.")
	}
}

func TestUserResource_updateUserPasswordResource_Failure(t *testing.T) {
	username := "betty"
	userSvc := new(mockservice.UserServiceMock)
	passwordReset := &model.PasswordReset{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}
	body := marshallModel(passwordReset)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	userSvc.On("UpdateUserPassword", username, passwordReset, false).Return(errors.New("Failure"))

	handler := updateUserPasswordResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusInternalServerError, err.(*smokeStatus).Code, "Invalid status.")
	}
}
func TestUserResource_updateUserPasswordResource_BadJson(t *testing.T) {
	username := "betty"
	userSvc := new(mockservice.UserServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader("Bad JSON"))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	handler := updateUserPasswordResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusBadRequest, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestUserResource_updateUserPasswordResource_NotAdminAndNotSelf(t *testing.T) {
	username := "betty"
	userSvc := new(mockservice.UserServiceMock)
	passwordReset := &model.PasswordReset{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}
	body := marshallModel(passwordReset)
	e := echo.New()
	req := test.NewRequest(echo.PUT, "/users/" + username, strings.NewReader(body))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)
	req.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "not" + username
	claims["isAdmin"] = false
	c.Set("user", token)
	c.SetParamNames("userid")
	c.SetParamValues(username)

	handler := updateUserPasswordResource(userSvc)

	err := handler(c)

	userSvc.AssertExpectations(t)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusForbidden, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func marshallModel(m interface{}) string {
	b, _ := json.Marshal(m)
	return string(b)
}

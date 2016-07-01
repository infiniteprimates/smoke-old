package server

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	mockserver "github.com/infiniteprimates/smoke/mocks/server"
	mockservice "github.com/infiniteprimates/smoke/mocks/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthResource_createAuthResources(t *testing.T) {
	router := new(mockserver.RouterMock)
	authSvc := new(mockservice.AuthServiceMock)

	router.On("POST", "/auth", mock.AnythingOfType("echo.HandlerFunc"), mock.AnythingOfType("[]echo.MiddlewareFunc"))
	createAuthResources(router, authSvc)
}

func TestAuthResource_postAuthorizationResource_success(t *testing.T) {
	username := "user"
	password := "pass"
	token := "asdfjkl;"
	authSvc := new(mockservice.AuthServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	c.Set("username", username)
	c.Set("password", password)

	authSvc.On("AuthenticateUser", username, password).Return(token, nil)

	handler := postAuthorizationResource(authSvc)
	err := handler(c)

	if assert.NoError(t, err, "An error occured invoking handler.") {
		assert.Equal(t, http.StatusOK, res.Status(), "Invalid status.")
		assert.Contains(t, res.Body.String(), token, "Invalid token.")
	}
}

func TestAuthResource_postAuthorizationResource_Unauthorized(t *testing.T) {
	username := "user"
	password := "pass"
	authSvc := new(mockservice.AuthServiceMock)
	e := echo.New()
	req := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	c.Set("username", username)
	c.Set("password", password)

	authSvc.On("AuthenticateUser", username, password).Return("", errors.New("failure"))

	handler := postAuthorizationResource(authSvc)
	err := handler(c)

	if assert.Error(t, err, "An error occured invoking handler.") {
		assert.Equal(t, http.StatusUnauthorized, err.(*smokeStatus).Code, "Invalid status.")
	}
}

func TestAuthResource_authorizationMiddleware(t *testing.T) {
	mw := authorizationMiddleware("jwtKey")
	assert.NotNil(t, mw, "Middleware was nil.")
}

func TestAuthResource_requireAdminMiddleware_success(t *testing.T) {
	successMsg := "Success"
	failureMsg := "Failure"
	e := echo.New()
	req := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["isAdmin"] = true

	c.Set("user", token)

	mw := requireAdminMiddleware(failureMsg)
	err := mw(func(ctx echo.Context) error { return c.String(http.StatusOK, successMsg) })(c)

	if assert.NoError(t, err, "An error occured invoking handler.") {
		assert.Equal(t, http.StatusOK, res.Status(), "Invalid status.")
		assert.Contains(t, res.Body.String(), successMsg, "Invalid body.")
	}
}

func TestAuthResource_requireAdminMiddleware_Forbidden(t *testing.T) {
	successMsg := "Success"
	failureMsg := "Failure"
	e := echo.New()
	req := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["isAdmin"] = false

	c.Set("user", token)

	mw := requireAdminMiddleware(failureMsg)
	err := mw(func(ctx echo.Context) error { return c.String(http.StatusOK, successMsg) })(c)

	if assert.Error(t, err, "Expected error not returned.") {
		assert.Equal(t, http.StatusForbidden, err.(*smokeStatus).Code, "Invalid status.")
		assert.Equal(t, err.(*smokeStatus).Message, failureMsg, "Invalid body.")
	}
}

func TestAuthResource_basicAuthExtractor_success(t *testing.T) {
	authHeader := "Basic YWRtaW46c2VjcmV0"
	e := echo.New()
	req := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header().Set("Authorization", authHeader)
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	mw := basicAuthExtractor()
	err := mw(func(ctx echo.Context) error { return c.String(http.StatusOK, "") })(c)

	if assert.NoError(t, err, "An error occured invoking handler.") {
		assert.Equal(t, "admin", c.Get("username"), "Invalid username.")
		assert.Equal(t, "secret", c.Get("password"), "Invalid password.")
	}
}

func TestAuthResource_basicAuthExtractor_bad(t *testing.T) {
	authHeader := "Bogus asdf"
	e := echo.New()
	req := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header().Set("Authorization", authHeader)
	res := test.NewResponseRecorder()
	c := e.NewContext(req, res)

	mw := basicAuthExtractor()
	err := mw(func(ctx echo.Context) error { return c.String(http.StatusOK, "") })(c)

	assert.Error(t, err, "Expected error not returned.")
}

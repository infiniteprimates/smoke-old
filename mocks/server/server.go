package server

import (
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

type (
	RouterMock struct {
		mock.Mock
	}
)

func (m *RouterMock) Any(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) File(path string, file string) {
	m.Called(path, file)
}

func (m *RouterMock) Group(path string, mw ...echo.MiddlewareFunc) *echo.Group {
	args := m.Called(path, mw)
	return args.Get(0).(*echo.Group)
}

func (m *RouterMock) Match(methods []string, path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(methods, path, h, mw)
}

func (m *RouterMock) Static(prefix string, root string) {
	m.Called(prefix, root)
}

func (m *RouterMock) Use(mw ...echo.MiddlewareFunc) {
	m.Called(mw)
}

func (m *RouterMock) CONNECT(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) DELETE(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) GET(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) HEAD(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) OPTIONS(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) PATCH(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) POST(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) PUT(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

func (m *RouterMock) TRACE(path string, h echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
	m.Called(path, h, mw)
}

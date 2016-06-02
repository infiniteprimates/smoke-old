package server

import (
	"fmt"
	"io"
	"net"

	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
)

type (
	server struct {
		*echo.Echo
		ip   string
		port uint16
	}

	router interface {
		Any(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		File(path string, file string)
		Group(path string, m ...echo.MiddlewareFunc) *echo.Group
		Match(methods []string, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		Static(prefix string, root string)
		Use(m ...echo.MiddlewareFunc)
		CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
		TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	}

	Server interface {
		Start()
	}
)

func New(logWriter io.Writer, cfg *config.Config, userService *service.UserService, authService *service.AuthService) (Server, error) {
	e := echo.New()

	e.SetHTTPErrorHandler(smokeErrorHandler(e))

	if debug := cfg.GetBool(config.Debug); debug {
		e.SetDebug(debug)
	}

	logConfig := middleware.DefaultLoggerConfig
	logConfig.Output = logWriter
	e.Use(middleware.LoggerWithConfig(logConfig))

	e.Use(middleware.Recover())

	if cfg.GetBool(config.DevCors) {
		corsConfig := middleware.DefaultCORSConfig
		corsConfig.AllowHeaders = []string{echo.HeaderOrigin, echo.HeaderAuthorization, echo.HeaderContentType}
		corsConfig.MaxAge = 60
		e.Use(middleware.CORSWithConfig(corsConfig))
	}

	root := cfg.GetString(config.UiRoot)
	e.Use(middleware.Static(root))

	createResources(e, cfg, userService, authService)

	ip := cfg.GetString(config.Ip)
	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("Configured listen ip '%s' is invalid", ip)
	}

	port := cfg.GetInt(config.Port)
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("Configured listen port '%s' is invalid", cfg.GetString(config.Port))
	}

	server := &server{
		Echo: e,
		ip:   ip,
		port: uint16(port),
	}

	return server, nil
}

func createResources(r router, cfg *config.Config, userService *service.UserService, authService *service.AuthService) {
	group := r.Group("/api")
	createAuthResources(group, authService)
	createUserResources(group, cfg, userService)
}

func (server *server) Start() {
	ipAndPort := fmt.Sprintf("%s:%d", server.ip, server.port)
	//TODO: Log ipAndPort here
	server.Run(fasthttp.New(ipAndPort))
}

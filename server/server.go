package server

import (
	"fmt"
	"io"
	"net"

	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

type Server interface {
	Start()
}

type server struct {
	*echo.Echo
	ip   string
	port uint16
}

func New(logWriter io.Writer, cfg *viper.Viper, db *db.Db) (Server, error) {
	e := echo.New()
	if(cfg.GetBool(config.DEBUG)) {
		e.SetDebug(false)
	}

	middleware.DefaultLoggerConfig.Output = logWriter
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	root := cfg.GetString(config.UI_ROOT)
	e.Use(middleware.Static(root))

	createResources(db, e)

	ip := cfg.GetString(config.IP)
	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("Configured listen ip '%s' is invalid", ip)
	}

	port := cfg.GetInt(config.PORT)
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("Configured listen port '%s' is invalid", cfg.GetString(config.PORT))
	}

	server := &server{
		Echo: e,
		ip:     ip,
		port:   uint16(port),
	}

	return server, nil
}

func createResources(db *db.Db, e *echo.Echo) {
	g := e.Group("/api")
	createAuthResources(db, g)
	createUserResources(db, g)
}

func (server *server) Start() {
	ipAndPort := fmt.Sprintf("%s:%d", server.ip, server.port)
	//TODO: Log ipAndPort here
	server.Run(fasthttp.New(ipAndPort))
}

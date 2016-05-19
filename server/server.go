package server

import (
	"fmt"
	"io"
	"net"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/db"
	mw "github.com/infiniteprimates/smoke/middleware"
	"github.com/spf13/viper"
)

type Server interface {
	Start()
}

type server struct {
	*gin.Engine
	ip   string
	port uint16
}

func New(logWriter io.Writer, cfg *viper.Viper, db *db.Db) (Server, error) {
	gin.SetMode(cfg.GetString(config.GIN_MODE))

	router := gin.New()
	router.Use(gin.LoggerWithWriter(logWriter))
	router.Use(gin.RecoveryWithWriter(logWriter))

	//TODO: Create a static content handler that works without directory listing.
	root := cfg.GetString(config.UI_ROOT)
	router.NoRoute(mw.MetricsHandler("static"), static.Serve("/", static.LocalFile(root, true)))

	createResources(db, router)

	ip := cfg.GetString(config.IP)
	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("Configured listen ip '%s' is invalid", ip)
	}

	port := cfg.GetInt(config.PORT)
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("Configured listen port '%s' is invalid", cfg.GetString(config.PORT))
	}

	server := &server{
		Engine: router,
		ip:     ip,
		port:   uint16(port),
	}

	return server, nil
}

func createResources(db *db.Db, router gin.IRouter) {
	createAuthResources(db, router)
	createUserResources(db, router)
}

func (server *server) Start() {
	ipAndPort := fmt.Sprintf("%s:%d", server.ip, server.port)
	//TODO: Log ipAndPort here
	server.Run(ipAndPort)
}

package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/db"
	"github.com/spf13/viper"
)

type Server interface {
	Start()
}

type server struct {
	*gin.Engine
}

func New(logWriter io.Writer, db *db.Db) Server {
	router := gin.New()
	router.Use(gin.LoggerWithWriter(logWriter))
	router.Use(gin.RecoveryWithWriter(logWriter))

	//TODO: Figure out how to do static better with router.NoRoute and contrib static. Doesn't work right though.
	router.Any("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusTemporaryRedirect, "/ui/") })
	router.Static("/ui", "ui")

	createResources(db, router)

	server := &server {
		Engine: router,
	}

	return server
}

func createResources(db *db.Db, router gin.IRouter) {
	createAuthResources(db, router)
	createUserResources(db, router)
}

func (server *server) Start() {
	ipAndPort := fmt.Sprintf("%s:%d", viper.GetString(config.IP), viper.GetInt(config.PORT))
	server.Run(ipAndPort)
}


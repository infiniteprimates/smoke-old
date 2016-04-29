package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/rest"
	"github.com/rcrowley/go-metrics"
	"github.com/Sirupsen/logrus"
)

func main() {

	//TODO: get some configuration going here
	startServer()
}

func startServer() {
	logger := logrus.New()
	logWriter := logger.Writer()
	defer logWriter.Close()

	// Start background metrics logger
	//TODO: make metrics logging interval configurable
	go metrics.Log(metrics.DefaultRegistry, 5 * time.Minute, log.New(logWriter, "metrics", log.Lmicroseconds))

	router := gin.New()
	router.Use(gin.LoggerWithWriter(logWriter))
	router.Use(gin.RecoveryWithWriter(logWriter))

	//TODO: Figure out how to do static better with router.NoRoute and contrib static. Doesn't work right though.
	router.Any("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusTemporaryRedirect, "/ui/")})
	router.Static("/ui", "ui")

	createResources(router)

	router.Run() //TODO: make port and listen address configurable
}

func createResources(router gin.IRouter) {
	rest.CreateAuthResources(router)
	rest.CreateUserResources(router)
}

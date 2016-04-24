package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/user"
	"github.com/rcrowley/go-metrics"
	"github.com/Sirupsen/logrus"
	"github.com/infiniteprimates/smoke/auth"
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

	createResources(router)

	router.Run() //TODO: make port and listen address configurable
}

func createResources(router gin.IRouter) {
	auth.CreateAuthResources(router)
	user.CreateUserResources(router)
}


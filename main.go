package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logWriter := logger.Writer()

	router := gin.New()
	router.Use(gin.LoggerWithWriter(logWriter))
	router.Use(gin.RecoveryWithWriter(logWriter))

	router.Run() //TODO: make configurable
}

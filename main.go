package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/rest"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
)

func main() {
	config.Init()
	startServer()
}

func startServer() {
	logger := logrus.New()
	logWriter := logger.Writer()
	defer logWriter.Close()

	// Start background metrics logger
	metricsLoggingInterval := time.Duration(viper.GetInt(config.METRICS_LOGGING_INTERVAL)) * time.Minute
	go metrics.Log(metrics.DefaultRegistry, metricsLoggingInterval, log.New(logWriter, "metrics", log.Lmicroseconds))

	router := gin.New()
	router.Use(gin.LoggerWithWriter(logWriter))
	router.Use(gin.RecoveryWithWriter(logWriter))

	//TODO: Figure out how to do static better with router.NoRoute and contrib static. Doesn't work right though.
	router.Any("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusTemporaryRedirect, "/ui/") })
	router.Static("/ui", "ui")

	createResources(router)

	ipAndPort := fmt.Sprintf("%s:%d", viper.GetString(config.IP), viper.GetInt(config.PORT))
	router.Run(ipAndPort)
}

func createResources(router gin.IRouter) {
	rest.CreateAuthResources(router)
	rest.CreateUserResources(router)
}

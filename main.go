package main

import (
	"io"
	"log"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/server"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"github.com/infiniteprimates/smoke/db"
)

func main() {
	logger := logrus.New()
	logWriter := logger.Writer()
	defer logWriter.Close()

	config.Init()

	startMetricsLogging(logWriter)

	db := db.New()

	server := server.New(logWriter, db)
	server.Start()
}

func startMetricsLogging(logWriter io.Writer) {
	// Start background metrics logger
	metricsLoggingInterval := time.Duration(viper.GetInt(config.METRICS_LOGGING_INTERVAL)) * time.Minute
	go metrics.Log(metrics.DefaultRegistry, metricsLoggingInterval, log.New(logWriter, "metrics", log.Lmicroseconds))
}

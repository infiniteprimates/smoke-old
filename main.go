package main

import (
	"io"
	"log"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/db"
	"github.com/infiniteprimates/smoke/server"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
)

func main() {
	logger := logrus.New()
	logWriter := logger.Writer()
	defer logWriter.Close()

	cfg := config.New()

	startMetricsLogging(cfg, logWriter)

	db := db.New()

	srv, err := server.New(logWriter, cfg, db)
	if err != nil {
		logger.Error(err)
		return
	}

	srv.Start()
}

func startMetricsLogging(cfg *viper.Viper, logWriter io.Writer) {
	// Start background metrics logger
	metricsLoggingInterval := time.Duration(cfg.GetInt(config.METRICS_LOGGING_INTERVAL)) * time.Minute
	go metrics.Log(metrics.DefaultRegistry, metricsLoggingInterval, log.New(logWriter, "metrics", log.Lmicroseconds))
}

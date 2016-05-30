package main

import (
	"io"
	"log"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/db"
	"github.com/infiniteprimates/smoke/model"
	"github.com/infiniteprimates/smoke/server"
	"github.com/infiniteprimates/smoke/service"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
)

func main() {
	logger := logrus.New()
	logWriter := logger.Writer()
	defer logWriter.Close()

	cfg := config.New()
	for k, v := range cfg.AllSettings() {
		logger.Infof("CONFIG: %s = %v", k, v)
	}

	startMetricsLogging(cfg, logWriter)

	userDb, err := db.NewUserDb(cfg)
	fatalIfErr(logger, err)

	passwordService := service.NewPasswordService()

	userService, err := service.NewUserService(userDb, passwordService)
	fatalIfErr(logger, err)
	initAccounts(userService) //TODO:temporary

	srv, err := server.New(logWriter, cfg, userService, passwordService)
	fatalIfErr(logger, err)

	srv.Start()
}

func fatalIfErr(logger *logrus.Logger, err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func startMetricsLogging(cfg *viper.Viper, logWriter io.Writer) {
	// Start background metrics logger
	metricsLoggingInterval := time.Duration(cfg.GetInt(config.MetricsLoggingInterval)) * time.Minute
	go metrics.Log(metrics.DefaultRegistry, metricsLoggingInterval, log.New(logWriter, "metrics", log.Lmicroseconds))
}

func initAccounts(userService *service.UserService) {
	// This is temporary code until we have a real DB
	userService.Create(&model.User{
		Username: "admin",
		Password: "secret",
		IsAdmin:  true,
	})
	userService.Create(&model.User{
		Username: "user",
		Password: "password",
		IsAdmin:  false,
	})
}

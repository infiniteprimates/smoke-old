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
)

func main() {
	//TODO: look into newer logging in echo.
	logger := logrus.New()
	logWriter := logger.Writer()
	defer logWriter.Close()

	cfg, err := config.GetConfig()
	fatalIfErr(logger, err)
	for k, v := range cfg.AllSettings() {
		logger.Infof("CONFIG: %s = %v", k, v)
	}

	startMetricsLogging(cfg, logWriter)

	userDb, err := db.NewUserDb(cfg)
	fatalIfErr(logger, err)

	authService := service.NewAuthService(cfg, userDb)

	userService := service.NewUserService(userDb, authService)

	//TODO:temporary account creation during initial dev
	initAccounts(userService)

	srv, err := server.New(logWriter, cfg, userService, authService)
	fatalIfErr(logger, err)

	srv.Start()
}

func fatalIfErr(logger *logrus.Logger, err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func startMetricsLogging(cfg config.Config, logWriter io.Writer) {
	// Start background metrics logger
	metricsLoggingInterval := time.Duration(cfg.GetInt(config.MetricsLoggingInterval)) * time.Minute
	go metrics.Log(metrics.DefaultRegistry, metricsLoggingInterval, log.New(logWriter, "metrics", log.Lmicroseconds))
}

func initAccounts(userService service.UserService) {
	// This is temporary code until we have a real DB and an admin bootstrapping process
	userService.Create(&model.User{
		Username: "admin",
		IsAdmin:  true,
	})
	userService.UpdateUserPassword("admin", &model.PasswordReset{NewPassword: "secret"}, true)

	userService.Create(&model.User{
		Username: "user",
		IsAdmin:  false,
	})
	userService.UpdateUserPassword("user", &model.PasswordReset{NewPassword: "password"}, true)
}

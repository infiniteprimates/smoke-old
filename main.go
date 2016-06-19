package main

import (
	"github.com/infiniteprimates/smoke/config"
	"github.com/infiniteprimates/smoke/db"
	"github.com/infiniteprimates/smoke/model"
	"github.com/infiniteprimates/smoke/server"
	"github.com/infiniteprimates/smoke/service"
	glog "github.com/labstack/gommon/log"
)

func main() {
	logger := glog.New("smoke")

	cfg, err := config.GetConfig()
	fatalIfErr(logger, err)
	for k, v := range cfg.AllSettings() {
		logger.Infof("CONFIG: %s = %v", k, v)
	}

	userDb, err := db.NewUserDb(cfg)
	fatalIfErr(logger, err)

	authService := service.NewAuthService(cfg, userDb)

	userService := service.NewUserService(userDb, authService)

	//TODO:temporary account creation during initial dev
	initAccounts(userService)

	srv, err := server.New(logger, cfg, userService, authService)
	fatalIfErr(logger, err)

	srv.Start()
}

func fatalIfErr(logger *glog.Logger, err error) {
	if err != nil {
		logger.Fatal(err)
	}
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

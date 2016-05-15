package config

import (
	"sync"

	"github.com/spf13/viper"
)

const (
	GIN_MODE                 = "gin_mode"
	IP                       = "ip"
	METRICS_LOGGING_INTERVAL = "metrics_logging_interval"
	PORT                     = "port"
	UI_ROOT                  = "ui_root"
)

var (
	config *viper.Viper
	once   sync.Once
)

func New() *viper.Viper {
	once.Do(func() {
		config = viper.New()

		config.AutomaticEnv()
		config.SetEnvPrefix("smoke")

		setDefaults(config)
	})

	return config
}

func setDefaults(config *viper.Viper) {
	config.SetDefault(GIN_MODE, "release")
	config.SetDefault(IP, "0.0.0.0")
	config.SetDefault(METRICS_LOGGING_INTERVAL, 5)
	config.SetDefault(PORT, 8080)
	config.SetDefault(UI_ROOT, "dist")
}

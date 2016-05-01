package config

import (
	"github.com/spf13/viper"
)

const (
	IP                       = "ip"
	METRICS_LOGGING_INTERVAL = "metrics_logging_interval"
	PORT                     = "port"
)

func Init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("smoke")

	setDefaults()
}

func setDefaults() {
	viper.SetDefault(IP, "0.0.0.0")
	viper.SetDefault(PORT, 8080)
	viper.SetDefault(METRICS_LOGGING_INTERVAL, 5)
}

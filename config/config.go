package config

import (
	"sync"

	"github.com/spf13/viper"
)

const (
	Debug                  = "debug"
	DevCors                = "dev_cors"
	Ip                     = "ip"
	MetricsLoggingInterval = "metrics_logging_interval"
	Port                   = "port"
	UiRoot                 = "ui_root"
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
	config.SetDefault(Debug, false)
	config.SetDefault(DevCors, false)
	config.SetDefault(Ip, "0.0.0.0")
	config.SetDefault(MetricsLoggingInterval, 5)
	config.SetDefault(Port, 8080)
	config.SetDefault(UiRoot, "dist")
}

package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		*viper.Viper
	}
)

const (
	Debug                  = "debug"
	DevCors                = "dev_cors"
	Ip                     = "ip"
	JwtKey                 = "jwt_key"
	MetricsLoggingInterval = "metrics_logging_interval"
	Port                   = "port"
	UiRoot                 = "ui_root"
)

var (
	config *Config
	once   sync.Once
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		viperConfig := viper.New()

		config = &Config{
			Viper: viperConfig,
		}

		config.AutomaticEnv()
		config.SetEnvPrefix("smoke")

		setDefaults(config)
	})

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults(config *Config) {
	config.SetDefault(Debug, false)
	config.SetDefault(DevCors, false)
	config.SetDefault(Ip, "0.0.0.0")
	config.SetDefault(MetricsLoggingInterval, 5)
	config.SetDefault(Port, 8080)
	config.SetDefault(UiRoot, "dist")
}

func validateConfig(config *Config) error {
	if !config.IsSet(JwtKey) {
		return fmt.Errorf("Configuration '%s' is required.", JwtKey)
	}

	return nil
}

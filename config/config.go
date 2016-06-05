package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		*viper.Viper
	}
)

const (
	EnvPrefix = "smoke"

	Debug                  = "debug"
	DevCors                = "dev_cors"
	Ip                     = "ip"
	JwtKey                 = "jwt_key"
	MetricsLoggingInterval = "metrics_logging_interval"
	Port                   = "port"
	UiRoot                 = "ui_root"

	defaultDebug                  = true
	defaultDevCors                = false
	defaultIp                     = "0.0.0.0"
	defaultMetricsLoggingInterval = 5
	defaultPort                   = 80
	defaultUiRoot                 = "dist"
)

func GetConfig() (*Config, error) {
	viperConfig := viper.New()

	config := &Config{
		Viper: viperConfig,
	}

	config.AutomaticEnv()
	config.SetEnvPrefix(EnvPrefix)

	setDefaults(config)

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults(config *Config) {
	config.SetDefault(Debug, defaultDebug)
	config.SetDefault(DevCors, defaultDevCors)
	config.SetDefault(Ip, defaultIp)
	config.SetDefault(MetricsLoggingInterval, defaultMetricsLoggingInterval)
	config.SetDefault(Port, defaultPort)
	config.SetDefault(UiRoot, defaultUiRoot)
}

func validateConfig(config *Config) error {
	var err error
	errMsg := ""

	if !config.IsSet(JwtKey) {
		errMsg += fmt.Sprintf("Configuration '%s' is required. ", JwtKey)
	}

	if len(errMsg) > 0 {
		err = fmt.Errorf(errMsg)
	}

	return err
}

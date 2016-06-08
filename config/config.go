package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config interface {
		AllSettings() map[string]interface{}
		GetBool(string) bool
		GetInt(string) int
		GetString(string) string
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
	defaultPort                   = 8080
	defaultUiRoot                 = "dist"
)

func GetConfig() (Config, error) {
	config := viper.New()

	config.AutomaticEnv()
	config.SetEnvPrefix(EnvPrefix)

	setDefaults(config)

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults(config *viper.Viper) {
	config.SetDefault(Debug, defaultDebug)
	config.SetDefault(DevCors, defaultDevCors)
	config.SetDefault(Ip, defaultIp)
	config.SetDefault(MetricsLoggingInterval, defaultMetricsLoggingInterval)
	config.SetDefault(Port, defaultPort)
	config.SetDefault(UiRoot, defaultUiRoot)
}

func validateConfig(config *viper.Viper) error {
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

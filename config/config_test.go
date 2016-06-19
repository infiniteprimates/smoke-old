package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	os.Clearenv()
	_, err := GetConfig()

	assert.Error(t, err, "Expected error not returned.")
}

func TestDefaults(t *testing.T) {
	os.Clearenv()
	initRequired()
	config, err := GetConfig()
	assert.Nil(t, err, "An error occured instantiating a config.")

	assert.Equal(t, defaultDebug, config.GetBool(Debug), "Config '%s' default is incorrect.", Debug)
	assert.Equal(t, defaultDevCors, config.GetBool(DevCors), "Config '%s' default is incorrect.", DevCors)
	assert.Equal(t, defaultIp, config.GetString(Ip), "Config '%s' default is incorrect.", Ip)
	assert.Equal(t, defaultMetricsPublishingInterval, config.GetInt(MetricsPublishingInterval), "Config '%s' default is incorrect.", MetricsPublishingInterval)
	assert.Equal(t, defaultPort, config.GetInt(Port), "Config '%s' default is incorrect.", Port)
	assert.Equal(t, defaultUiRoot, config.GetString(UiRoot), "Config '%s' default is incorrect.", UiRoot)
}

func TestValues(t *testing.T) {
	os.Clearenv()
	initRequired()
	config, err := GetConfig()
	assert.Nil(t, err, "An error occured instantiating a config.")

	ip := "ip"
	metricsLoggingInterval := 120
	port := 65535
	uiRoot := "uiRoot"

	os.Setenv(createEnvName(Debug), strconv.FormatBool(!defaultDebug))
	os.Setenv(createEnvName(DevCors), strconv.FormatBool(!defaultDevCors))
	os.Setenv(createEnvName(Ip), ip)
	os.Setenv(createEnvName(MetricsPublishingInterval), strconv.Itoa(metricsLoggingInterval))
	os.Setenv(createEnvName(Port), strconv.Itoa(port))
	os.Setenv(createEnvName(UiRoot), uiRoot)

	assert.Equal(t, !defaultDebug, config.GetBool(Debug), "Config '%s' default is incorrect.", Debug)
	assert.Equal(t, !defaultDevCors, config.GetBool(DevCors), "Config '%s' default is incorrect.", DevCors)
	assert.Equal(t, ip, config.GetString(Ip), "Config '%s' default is incorrect.", Ip)
	assert.Equal(t, metricsLoggingInterval, config.GetInt(MetricsPublishingInterval), "Config '%s' default is incorrect.", MetricsPublishingInterval)
	assert.Equal(t, port, config.GetInt(Port), "Config '%s' default is incorrect.", Port)
	assert.Equal(t, uiRoot, config.GetString(UiRoot), "Config '%s' default is incorrect.", UiRoot)

}

func initRequired() {
	os.Setenv(createEnvName(JwtKey), JwtKey)
}

func createEnvName(configName string) string {
	return strings.ToUpper(fmt.Sprintf("%s_%s", EnvPrefix, configName))
}

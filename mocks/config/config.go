package config

import (
	"github.com/stretchr/testify/mock"
)

type (
	ConfigMock struct {
		mock.Mock
	}
)

func (m *ConfigMock) AllSettings() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

func (m *ConfigMock) GetBool(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *ConfigMock) GetInt(name string) int {
	args := m.Called(name)
	return args.Int(0)
}

func (m *ConfigMock) GetString(name string) string {
	args := m.Called(name)
	return args.String(0)
}

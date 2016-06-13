package service

import "github.com/stretchr/testify/mock"

type(
	AuthServiceMock struct {
		mock.Mock
	}
)

func (m *AuthServiceMock) AuthenticateUser(username string, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *AuthServiceMock) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

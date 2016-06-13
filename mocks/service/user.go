package service

import (
	"github.com/stretchr/testify/mock"
	"github.com/infiniteprimates/smoke/model"
)

type (
	UserServiceMock struct {
		mock.Mock
	}
)

func (m *UserServiceMock) Create(userModel *model.User) (*model.User, error) {
	args := m.Called(userModel)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}

func (m *UserServiceMock) Find(username string) (*model.User, error) {
	args := m.Called(username)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}

func (m *UserServiceMock) List() ([]*model.User, error) {
	args := m.Called()
	users, _ := args.Get(0).([]*model.User)
	return users, args.Error(1)
}

func (m *UserServiceMock) Update(userModel *model.User) (*model.User, error) {
	args := m.Called(userModel)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}

func (m *UserServiceMock) Delete(username string) error {
	args := m.Called(username)
	return args.Error(0)

}

func (m *UserServiceMock) UpdateUserPassword(username string, passwordReset *model.PasswordReset, administrativeReset bool) error {
	args := m.Called(username, passwordReset, administrativeReset)
	return args.Error(0)
}

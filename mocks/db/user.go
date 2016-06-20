package db

import (
	"github.com/infiniteprimates/smoke/db"
	"github.com/stretchr/testify/mock"
)

type (
	UserDbMock struct {
		mock.Mock
	}
)

func (m *UserDbMock) Create(user *db.User) (*db.User, error) {
	args := m.Called(user)
	userRet, _ := args.Get(0).(*db.User)
	return userRet, args.Error(1)
}

func (m *UserDbMock) Find(username string) (*db.User, error) {
	args := m.Called(username)
	userRet, _ := args.Get(0).(*db.User)
	return userRet, args.Error(1)
}

func (m *UserDbMock) List() ([]*db.User, error) {
	args := m.Called()
	usersRet, _ := args.Get(0).([]*db.User)
	return usersRet, args.Error(1)
}

func (m *UserDbMock) Update(user *db.User) (*db.User, error) {
	args := m.Called(user)
	userRet, _ := args.Get(0).(*db.User)
	return userRet, args.Error(1)
}

func (m *UserDbMock) Delete(username string) error {
	args := m.Called(username)
	return args.Error(0)
}

func (m *UserDbMock) UpdateUserPassword(username string, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

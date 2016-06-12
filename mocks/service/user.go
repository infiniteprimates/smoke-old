package service

import "github.com/stretchr/testify/mock"

type (
	UserServiceMock struct {
		mock.Mock
	}
)

func (m *UserServiceMock) Create(userModel *model.User) (*model.User, error)
Find(username string) (*model.User, error)
List() ([]*model.User, error)
Update(userModel *model.User) (*model.User, error)
Delete(username string) error
UpdateUserPassword(username string, passwordReset *model.PasswordReset, administrativeReset bool) error

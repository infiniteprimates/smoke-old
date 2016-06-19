package service

import (
	"errors"
	"strings"

	"github.com/infiniteprimates/smoke/db"
	"github.com/infiniteprimates/smoke/model"
)

type (
	UserService interface {
		Create(userModel *model.User) (*model.User, error)
		Find(username string) (*model.User, error)
		List() ([]*model.User, error)
		Update(userModel *model.User) (*model.User, error)
		Delete(username string) error
		UpdateUserPassword(username string, passwordReset *model.PasswordReset, administrativeReset bool) error
	}

	userService struct {
		userDb      db.UserDb
		authService AuthService
	}
)

func NewUserService(userDb db.UserDb, authSvc AuthService) UserService {
	return &userService{
		userDb:      userDb,
		authService: authSvc,
	}
}

func (s *userService) Create(userModel *model.User) (*model.User, error) {
	userEntity := s.userModelToEntity(userModel)

	userEntity, err := s.userDb.Create(userEntity)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity), nil
}

func (s *userService) Find(userId string) (*model.User, error) {
	userEntity, err := s.userDb.Find(userId)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity), nil
}

func (s *userService) List() ([]*model.User, error) {
	userEntityList, err := s.userDb.List()
	if err != nil {
		return nil, err
	}

	userList := make([]*model.User, len(userEntityList), len(userEntityList))

	for i, userEntity := range userEntityList {
		userList[i] = s.userEntityToModel(userEntity)
	}

	return userList, nil
}

func (s *userService) Update(userModel *model.User) (*model.User, error) {
	userEntity := s.userModelToEntity(userModel)

	userEntity, err := s.userDb.Update(userEntity)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity), nil
}

func (s *userService) Delete(userId string) error {
	return s.userDb.Delete(userId)
}

func (s *userService) UpdateUserPassword(userId string, passwordReset *model.PasswordReset, administrativeReset bool) error {
	if len(strings.TrimSpace(passwordReset.NewPassword)) == 0 {
		return errors.New("Empty password is not acceptable.")
	}

	if !administrativeReset {
		if _, err := s.authService.AuthenticateUser(userId, passwordReset.OldPassword); err != nil {
			return err
		}
	}

	hashedPassword, err := s.authService.HashPassword(passwordReset.NewPassword)
	if err != nil {
		return err
	}

	return s.userDb.UpdateUserPassword(userId, hashedPassword)
}

func (s *userService) userEntityToModel(userEntity *db.User) *model.User {
	userModel := &model.User{
		Username: userEntity.Username,
		IsAdmin:  userEntity.IsAdmin,
	}

	return userModel
}

func (s *userService) userModelToEntity(userModel *model.User) *db.User {
	userEntity := &db.User{
		Username: userModel.Username,
		IsAdmin:  userModel.IsAdmin,
	}

	return userEntity
}

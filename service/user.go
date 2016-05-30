package service

import (
	"errors"
	"strings"

	"github.com/infiniteprimates/smoke/db"
	"github.com/infiniteprimates/smoke/model"
)

type UserService struct {
	userDb          *db.UserDb
	passwordService *PasswordService
}

func NewUserService(userDb *db.UserDb, passwordService *PasswordService) (*UserService, error) {
	return &UserService{
		passwordService: passwordService,
	}, nil
}

func (s *UserService) Create(userModel *model.User) (*model.User, error) {
	if len(strings.TrimSpace(userModel.Password)) == 0 {
		return nil, errors.New("User password is not acceptable")
	}

	userEntity, err := s.userModelToEntity(userModel)
	if err != nil {
		return nil, err
	}

	userEntity, err = s.userDb.Create(userEntity)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity, false), nil
}

func (s *UserService) Find(userId string, withPassword bool) (*model.User, error) {
	userEntity, err := s.userDb.Find(userId)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity, withPassword), nil
}

func (s *UserService) List() ([]*model.User, error) {
	userEntityList, err := s.userDb.List()
	if err != nil {
		return nil, err
	}

	userList := make([]*model.User, len(userEntityList), len(userEntityList))

	for i, userEntity := range userEntityList {
		userList[i] = s.userEntityToModel(userEntity, false)
	}

	return userList, nil
}

func (s *UserService) Update(userModel *model.User) (*model.User, error) {
	userEntity, err := s.userDb.Find(userModel.Username)
	if err != nil {
		return nil, err
	}

	userEntityFromModel, err := s.userModelToEntity(userModel)
	if err != nil {
		return nil, err
	}

	userEntity.IsAdmin = userEntityFromModel.IsAdmin
	if len(strings.TrimSpace(userEntityFromModel.Password)) > 0 {
		userEntity.Password = userEntityFromModel.Password
	}

	userEntity, err = s.userDb.Update(userEntity)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity, false), nil
}

func (s *UserService) Delete(userId string) error {
	return s.userDb.Delete(userId)
}

func (s *UserService) userEntityToModel(userEntity *db.User, withPassword bool) *model.User {
	userModel := &model.User{
		Username: userEntity.Username,
		IsAdmin:  userEntity.IsAdmin,
	}

	if withPassword {
		userModel.Password = userEntity.Password
	}

	return userModel
}

func (s *UserService) userModelToEntity(userModel *model.User) (*db.User, error) {
	userEntity := &db.User{
		Username: userModel.Username,
		IsAdmin:  userModel.IsAdmin,
	}

	if len(strings.TrimSpace(userModel.Password)) > 0 {
		hashedPassword, err := s.passwordService.HashPassword(userModel.Password)
		if err != nil {
			return nil, err
		}

		userEntity.Password = hashedPassword
	}

	return userEntity, nil
}

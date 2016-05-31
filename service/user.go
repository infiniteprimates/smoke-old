package service

import (
	"errors"
	"strings"

	"github.com/infiniteprimates/smoke/db"
	"github.com/infiniteprimates/smoke/model"
)

type UserService struct {
	userDb          *db.UserDb
	authService *AuthService
}

//TODO: WTH userDb!?!?!
func NewUserService(userDb *db.UserDb, authService *AuthService) (*UserService, error) {
	return &UserService{
		authService: authService,
	}, nil
}

func (s *UserService) Create(userModel *model.User) (*model.User, error) {
	userEntity := s.userModelToEntity(userModel)

	userEntity, err := s.userDb.Create(userEntity)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity), nil
}

func (s *UserService) Find(userId string) (*model.User, error) {
	userEntity, err := s.userDb.Find(userId)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity), nil
}

func (s *UserService) List() ([]*model.User, error) {
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

func (s *UserService) Update(userModel *model.User) (*model.User, error) {
	userEntity := s.userModelToEntity(userModel)

	userEntity, err := s.userDb.Update(userEntity)
	if err != nil {
		return nil, err
	}

	return s.userEntityToModel(userEntity), nil
}

func (s *UserService) Delete(userId string) error {
	return s.userDb.Delete(userId)
}

func (s *UserService) UpdateUserPassword(userId string, password string) error {
	if len(strings.TrimSpace(password)) == 0 {
		return errors.New("Empty password is not acceptable.")
	}

	hashedPassword, err := s.authService.hashPassword(password)
	if err != nil {
		return err
	}

	return s.userDb.UpdateUserPassword(userId, hashedPassword)
}

func (s *UserService) userEntityToModel(userEntity *db.User) *model.User {
	userModel := &model.User{
		Username: userEntity.Username,
		IsAdmin:  userEntity.IsAdmin,
	}

	return userModel
}

func (s *UserService) userModelToEntity(userModel *model.User) *db.User {
	userEntity := &db.User{
		Username: userModel.Username,
		IsAdmin:  userModel.IsAdmin,
	}

	return userEntity
}

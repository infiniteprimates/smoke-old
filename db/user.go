package db

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type (
	UserDb struct{}

	User struct {
		Username string
		Password string
		IsAdmin  bool
	}
)

// Yes, I know the map isn't threadsafe. It's temporary.
var users = map[string]User{}

func NewUserDb(cfg *viper.Viper) (*UserDb, error) {
	return &UserDb{}, nil
}

func (db *UserDb) Create(user *User) (*User, error) {
	if _, present := users[user.Username]; present {
		return nil, errors.New("User already exists")
	}

	users[user.Username] = *user

	return user, nil
}

func (db *UserDb) Find(userId string) (*User, error) {
	userEntity, present := users[userId]
	if !present {
		return nil, NewDbError(DbNotFound, fmt.Sprintf("User '%s' not found."))
	}

	return &userEntity, nil
}

func (db *UserDb) List() ([]*User, error) {
	userList := make([]*User, len(users), len(users))
	i := 0
	for _, user := range users {
		//TODO:Temp can be removed when this is a real database
		temp := user
		userList[i] = &temp
		i++
	}
	return userList, nil
}

func (db *UserDb) Update(user *User) (*User, error) {
	users[user.Username] = *user
	return user, nil
}

func (db *UserDb) Delete(userId string) error {
	delete(users, userId)
	return nil
}

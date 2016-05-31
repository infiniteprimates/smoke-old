package db

import (
	"github.com/infiniteprimates/smoke/config"
)

type (
	UserDb struct{}

	User struct {
		Username string
		PasswordHash string
		IsAdmin  bool
	}
)

//TODO: Yes, I know the map isn't threadsafe. It's temporary.
var users = map[string]User{}

func NewUserDb(cfg *config.Config) (*UserDb, error) {
	return &UserDb{}, nil
}

func (db *UserDb) Create(user *User) (*User, error) {
	if _, present := users[user.Username]; present {
		return nil, NewDbError(EntityExists, "User '%s' already exists", user.Username)
	}

	users[user.Username] = *user

	return user, nil
}

func (db *UserDb) Find(userId string) (*User, error) {
	userEntity, present := users[userId]
	if !present {
		return nil, NewDbError(EntityNotFound, "User '%s' not found.", userId)
	}

	return &userEntity, nil
}

func (db *UserDb) List() ([]*User, error) {
	userList := make([]*User, len(users), len(users))
	i := 0
	for _, user := range users {
		//TODO: Temp can be removed when this is a real database. It's needed to force a copy for now.
		temp := user
		userList[i] = &temp
		i++
	}
	return userList, nil
}

func (db *UserDb) Update(user *User) (*User, error) {
	//TODO: This is a poor mans update keeping the password hash intact.
	if origUser, present := users[user.Username]; !present {
		return nil, NewDbError(EntityNotFound, "User '%s' not found.", user.Username)
	} else {
		user.PasswordHash = origUser.PasswordHash
	}

	users[user.Username] = *user
	return user, nil
}

func (db *UserDb) Delete(userId string) error {
	delete(users, userId)
	return nil
}

func (db *UserDb) UpdateUserPassword(userId string, passwordhash string) error {
	if user, present := users[userId]; !present {
		return NewDbError(EntityNotFound, "User '%s' not found.", userId)
	} else {
		user.PasswordHash = passwordhash
		users[userId] = user
	}

	return nil
}

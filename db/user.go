package db

import (
	"github.com/infiniteprimates/smoke/config"
)

type (
	User struct {
		Username     string
		PasswordHash string
		IsAdmin      bool
	}

	UserDb interface {
		Create(user *User) (*User, error)
		Find(username string) (*User, error)
		List() ([]*User, error)
		Update(user *User) (*User, error)
		Delete(username string) error
		UpdateUserPassword(username string, password string) error
	}

	userDb struct {
		users map[string]User
	}
)

//TODO: Yes, I know the map isn't threadsafe. It's temporary.
var users map[string]User = map[string]User{}

func NewUserDb(cfg config.Config) (UserDb, error) {
	return &userDb{
		users: users,
	}, nil
}

func (db *userDb) Create(user *User) (*User, error) {
	if _, present := db.users[user.Username]; present {
		return nil, newDbError(EntityExists, "User '%s' already exists", user.Username)
	}

	db.users[user.Username] = *user

	return user, nil
}

func (db *userDb) Find(userId string) (*User, error) {
	userEntity, present := db.users[userId]
	if !present {
		return nil, newDbError(EntityNotFound, "User '%s' not found.", userId)
	}

	return &userEntity, nil
}

func (db *userDb) List() ([]*User, error) {
	userList := make([]*User, len(db.users), len(db.users))
	i := 0
	for _, user := range db.users {
		//TODO: Temp can be removed when this is a real database. It's needed to force a copy for now.
		temp := user
		userList[i] = &temp
		i++
	}
	return userList, nil
}

func (db *userDb) Update(user *User) (*User, error) {
	//TODO: This is a poor mans update keeping the password hash intact.
	origUser, present := db.users[user.Username]
	if !present {
		return nil, newDbError(EntityNotFound, "User '%s' not found.", user.Username)
	}

	user.PasswordHash = origUser.PasswordHash

	db.users[user.Username] = *user
	return user, nil
}

func (db *userDb) Delete(userId string) error {
	delete(db.users, userId)
	return nil
}

func (db *userDb) UpdateUserPassword(userId string, passwordhash string) error {
	user, present := db.users[userId]
	if !present {
		return newDbError(EntityNotFound, "User '%s' not found.", userId)
	}

	user.PasswordHash = passwordhash
	db.users[userId] = user

	return nil
}

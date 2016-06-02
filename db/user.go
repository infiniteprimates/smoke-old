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

	UserDb struct {
		//TODO: Yes, I know the map isn't threadsafe. It's temporary.
		users map[string]User
	}
)

func NewUserDb(cfg *config.Config) (*UserDb, error) {
	return &UserDb{
		users: map[string]User{},
	}, nil
}

func (db *UserDb) Create(user *User) (*User, error) {
	if _, present := db.users[user.Username]; present {
		return nil, NewDbError(EntityExists, "User '%s' already exists", user.Username)
	}

	db.users[user.Username] = *user

	return user, nil
}

func (db *UserDb) Find(userId string) (*User, error) {
	userEntity, present := db.users[userId]
	if !present {
		return nil, NewDbError(EntityNotFound, "User '%s' not found.", userId)
	}

	return &userEntity, nil
}

func (db *UserDb) List() ([]*User, error) {
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

func (db *UserDb) Update(user *User) (*User, error) {
	//TODO: This is a poor mans update keeping the password hash intact.
	origUser, present := db.users[user.Username]
	if !present {
		return nil, NewDbError(EntityNotFound, "User '%s' not found.", user.Username)
	}

	user.PasswordHash = origUser.PasswordHash

	db.users[user.Username] = *user
	return user, nil
}

func (db *UserDb) Delete(userId string) error {
	delete(db.users, userId)
	return nil
}

func (db *UserDb) UpdateUserPassword(userId string, passwordhash string) error {
	user, present := db.users[userId]
	if !present {
		return NewDbError(EntityNotFound, "User '%s' not found.", userId)
	}

	user.PasswordHash = passwordhash
	db.users[userId] = user

	return nil
}

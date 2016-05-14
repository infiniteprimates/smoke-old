package db

import (
	"github.com/infiniteprimates/smoke/model"
	"golang.org/x/crypto/bcrypt"
)

type user struct{}

var users = map[string]model.User{}

func init() {
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	users["admin"] = model.User{
		Username: "admin",
		Password: string(adminPassword),
		IsAdmin:  true,
	}

	userPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	users["user"] = model.User{
		Username: "user",
		Password: string(userPassword),
		IsAdmin:  false,
	}
}

func (user *user) FindUser(username string) (*model.User, error) {
	if user, present := users[username]; !present {
		return nil, nil
	} else {
		return &user, nil
	}
}

func (user *user) ListUsers() ([]model.User, error) {
	result := make([]model.User, 0, len(users))
	for _, user := range users {
		result = append(result, user)
	}
	return result, nil
}

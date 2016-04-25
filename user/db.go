package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var users = map[string]User{}

func init() {
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	users["admin"] = User {
		Username: "admin",
		Password: string(adminPassword),
		IsAdmin: true,
	}
}

func Find(username string) (*User, error) {
	if user, present := users[username] ; !present {
		return nil, errors.New("Not found")
	} else {
		return &user, nil
	}
}

func List() ([]User, error) {
	result := make([]User, 0, len(users))
	for _,user := range users {
		result = append(result, user)
	}
	return result, nil
}
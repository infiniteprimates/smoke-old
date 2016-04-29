package user

import (
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

	userPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	users["user"] = User {
		Username: "user",
		Password: string(userPassword),
		IsAdmin: false,
	}
}

func Find(username string) (*User, error) {
	if user, present := users[username] ; !present {
		return nil, nil
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
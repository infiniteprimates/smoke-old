package model

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"isAdmin"`
}
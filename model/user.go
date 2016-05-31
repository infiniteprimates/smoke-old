package model

type User struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

package model

type (
	Auth struct {
		AuthType string `json:"type"`
		Token    string `json:"token"`
	}
)

package model

type (
	User struct {
		Username string `json:"username"`
		IsAdmin  bool   `json:"isAdmin"`
	}

	PasswordReset struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
)

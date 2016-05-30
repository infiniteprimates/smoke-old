package service

import "golang.org/x/crypto/bcrypt"

type (
	PasswordService struct{}
)

func NewPasswordService() *PasswordService {
	return new(PasswordService)
}

func (s *PasswordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (s *PasswordService) ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

package service

import (
	"testing"

	"github.com/infiniteprimates/smoke/db"
	mockdb "github.com/infiniteprimates/smoke/mocks/db"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create_Success(b *testing.B) {
	userDb := new(mockdb.UserDbMock)
	authService := new(mockservice.UserServiceMock)
	svc := NewAuthService(userDb, authSvc)
}
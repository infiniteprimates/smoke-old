package db

import (
	"testing"

	"github.com/infiniteprimates/smoke/mocks/config"
	"github.com/stretchr/testify/assert"
)

func TestUserDb_Create_Success(t *testing.T) {
	cfg := new(config.ConfigMock)
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{}
		result, err := db.Create(&user)
		if assert.NoError(t, err, "An error occured creating a user.") {
			assert.Equal(t, &user, result, "Users mismatch.")
		}
	}
}

func TestUserDb_Create_AlreadyExists(t *testing.T) {
	cfg := new(config.ConfigMock)
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{ user.Username : user }
		_, err := db.Create(&user)
		if assert.Error(t, err, "Expected error not returned.") {
			if dbErr, ok := err.(*dbError); ok {
				assert.Equal(t, dbErrorReason(EntityExists), dbErr.reason, "Error code was not EntityExists.")
			} else {
				assert.Fail(t, "Returned error was not a dbError.")
			}
		}
	}
}

func TestUserDb_Find_Success(t *testing.T) {
	cfg := new(config.ConfigMock)
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{user.Username: user}
		result, err := db.Find(user.Username)
		if assert.NoError(t, err, "An error occured retrieving a user.") {
			assert.Equal(t, &user, result, "Expected user not returned.")
		}
	}
}

func TestUserDb_Find_NotFound(t *testing.T) {
	cfg := new(config.ConfigMock)
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{user.Username : user}
		_, err := db.Find("bogus")
		if assert.Error(t, err, "Expected error not returned.") {
			if dbErr, ok := err.(*dbError); ok {
				assert.Equal(t, dbErrorReason(EntityNotFound), dbErr.reason, "Error code was not EntityNotFound.")
			} else {
				assert.Fail(t, "Returned error was not a dbError.")
			}
		}
	}
}

func TestUserDb_List_Success(t *testing.T) {
	cfg := new(config.ConfigMock)
	admin := User {
		Username: "admin",
	}
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{ admin.Username : admin, user.Username: user }
		result, err := db.List()
		if assert.NoError(t, err, "An error occured while listing users.") && assert.NotNil(t, result, "List result was nil.") {
			assert.Len(t, result, 2, "List size is incorrect.")
			assert.Contains(t, result, &admin, "User list does not contain expected element.")
			assert.Contains(t, result, &user, "User list does not contain expected element.")
		}
	}
}

func TestUserDb_Update_Exists(t *testing.T) {
	cfg := new(config.ConfigMock)
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{ user.Username: user }
		result, err := db.Update(&user)
		if assert.NoError(t, err, "An error occured creating a user.") {
			assert.Equal(t, &user, result, "Users mismatch.")
		}
	}
}

func TestUserDb_Update_NotExists(t *testing.T) {
	cfg := new(config.ConfigMock)
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{}
		_, err := db.Update(&user)
		if assert.Error(t, err, "Expected error not returned.") {
			if dbErr, ok := err.(*dbError); ok {
				assert.Equal(t, dbErrorReason(EntityNotFound), dbErr.reason, "Error code was not EntityNotFound.")
			} else {
				assert.Fail(t, "Returned error was not a dbError.")
			}
		}
	}

}

func TestUserDb_Delete_Exists(t *testing.T) {
	cfg := new(config.ConfigMock)
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	db.(*userDb).users = map[string]User{ user.Username: user}
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		err := db.Delete(user.Username)
		assert.NoError(t, err, "An error occured deleting a user.")
	}
}

func TestUserDb_Delete_NotExists(t *testing.T) {
	cfg := new(config.ConfigMock)
	db, err := NewUserDb(cfg)
	db.(*userDb).users = map[string]User{}
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		err := db.Delete("bogus")
		assert.NoError(t, err, "An error occured deleting a user.")
	}
}

func TestUserDb_UpdateUserPassword_Exists(t *testing.T) {
	cfg := new(config.ConfigMock)
	user := User {
		Username: "user",
	}
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{ user.Username: user }
		err := db.UpdateUserPassword(user.Username, "hash")
		assert.NoError(t, err, "An error occured updating a user's password.")
	}
}

func TestUserDb_UpdateUserPassword_NotExists(t *testing.T) {
	cfg := new(config.ConfigMock)
	db, err := NewUserDb(cfg)
	if assert.NoError(t, err, "An error occured instantiating a UserDb.") {
		db.(*userDb).users = map[string]User{}
		err := db.UpdateUserPassword("bogus", "hash")
		if assert.Error(t, err, "Expected error not returned.") {
			if dbErr, ok := err.(*dbError); ok {
				assert.Equal(t, dbErrorReason(EntityNotFound), dbErr.reason, "Error code was not EntityNotFound.")
			} else {
				assert.Fail(t, "Returned error was not a dbError.")
			}
		}
	}

}

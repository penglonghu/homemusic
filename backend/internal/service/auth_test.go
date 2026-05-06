package service

import (
	"testing"

	"github.com/homemusic/backend/config"
	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthServiceTestDB(t *testing.T) {
	config.Conf = &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expire: 7200,
		},
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.User{})
	assert.NoError(t, err)

	dao.DB = db
}

func TestLoginSuccess(t *testing.T) {
	setupAuthServiceTestDB(t)

	user := &model.User{Username: "adminuser", Password: "Password123"}
	err := dao.CreateUser(user)
	assert.NoError(t, err)

	result, err := Login("adminuser", "Password123", false, "127.0.0.1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, "adminuser", result.Username)
	assert.Equal(t, config.Conf.JWT.Expire, result.Expires)
}

func TestLoginRememberMeExpire(t *testing.T) {
	setupAuthServiceTestDB(t)

	user := &model.User{Username: "rememberme", Password: "Password123"}
	err := dao.CreateUser(user)
	assert.NoError(t, err)

	result, err := Login("rememberme", "Password123", true, "127.0.0.2")
	assert.NoError(t, err)
	assert.Equal(t, rememberExpireSeconds, result.Expires)
}

func TestLoginFailedAndLockIP(t *testing.T) {
	setupAuthServiceTestDB(t)

	user := &model.User{Username: "lockuser", Password: "Password123"}
	err := dao.CreateUser(user)
	assert.NoError(t, err)

	for i := 0; i < maxFailedAttempts; i++ {
		_, err := Login("lockuser", "WrongPass123", false, "127.0.0.3")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	}

	_, err = Login("lockuser", "WrongPass123", false, "127.0.0.3")
	assert.ErrorIs(t, err, ErrIPLocked)
}

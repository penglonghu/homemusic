package service

import (
	"os"
	"testing"

	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupServiceTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&model.User{}, &model.SysInit{})
	assert.NoError(t, err)
	dao.DB = db
	dao.NewInitDao.db = db
}

func TestCheckInitStatus(t *testing.T) {
	setupServiceTestDB(t)

	status, err := CheckInitStatus()
	assert.NoError(t, err)
	assert.False(t, status["is_initialized"].(bool))
}

func TestExecuteInitialization(t *testing.T) {
	setupServiceTestDB(t)

	tempDir := t.TempDir()
	err := ExecuteInitialization("admin", "password123", "admin@example.com", tempDir)
	assert.NoError(t, err)

	status, err := CheckInitStatus()
	assert.NoError(t, err)
	assert.True(t, status["is_initialized"].(bool))
	assert.Equal(t, tempDir, status["music_path"].(string))
}

func TestExecuteInitializationInvalidInput(t *testing.T) {
	setupServiceTestDB(t)

	err := ExecuteInitialization("ad", "123", "", "/tmp")
	assert.Error(t, err)
}

func TestExecuteInitializationAlreadyInitialized(t *testing.T) {
	setupServiceTestDB(t)
	tempDir := t.TempDir()
	err := ExecuteInitialization("admin", "password123", "admin@example.com", tempDir)
	assert.NoError(t, err)

	err = ExecuteInitialization("admin2", "password123", "admin2@example.com", tempDir)
	assert.Error(t, err)
}

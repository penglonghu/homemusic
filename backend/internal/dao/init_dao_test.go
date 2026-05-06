package dao

import (
	"testing"
	"time"

	"github.com/homemusic/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDAOTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&model.SysInit{})
	assert.NoError(t, err)
	DB = db
	NewInitDao.db = db
}

func TestGetInitStatus(t *testing.T) {
	setupDAOTestDB(t)

	status, err := NewInitDao.GetInitStatus()
	assert.NoError(t, err)
	assert.False(t, status.IsInit)
}

func TestCreateAndUpdateInitConfig(t *testing.T) {
	setupDAOTestDB(t)

	config := &model.SysInit{
		IsInit:    true,
		MusicPath: "/music",
	}
	assert.NoError(t, NewInitDao.CreateInitConfig(config))
	assert.NotZero(t, config.ID)

	config.MusicPath = "/music/updated"
	assert.NoError(t, NewInitDao.UpdateInitConfig(config))

	updated, err := NewInitDao.GetInitStatus()
	assert.NoError(t, err)
	assert.Equal(t, "/music/updated", updated.MusicPath)
}

func TestIsInitialized(t *testing.T) {
	setupDAOTestDB(t)

	initialized, err := NewInitDao.IsInitialized()
	assert.NoError(t, err)
	assert.False(t, initialized)

	config := &model.SysInit{
		IsInit:    true,
		MusicPath: "/music",
	}
	assert.NoError(t, NewInitDao.CreateInitConfig(config))

	initialized, err = NewInitDao.IsInitialized()
	assert.NoError(t, err)
	assert.True(t, initialized)
}

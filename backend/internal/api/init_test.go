package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAPIInitTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.User{}, &model.SysInit{})
	assert.NoError(t, err)

	dao.DB = db
	dao.NewInitDao.db = db
}

func TestCheckInit(t *testing.T) {
	setupAPIInitTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/api/v1/init/check", CheckInit)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/init/check", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, false, resp["data"].(map[string]interface{})["is_initialized"])
}

func TestExecInit(t *testing.T) {
	setupAPIInitTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/api/v1/init/exec", ExecInit)

	tempDir := t.TempDir()
	payload := map[string]interface{}{
		"username":   "admin",
		"password":   "password123",
		"email":      "admin@example.com",
		"music_path": tempDir,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/init/exec", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "初始化完成", resp["data"].(map[string]interface{})["message"])
}

func TestExecInitInvalidInput(t *testing.T) {
	setupAPIInitTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/api/v1/init/exec", ExecInit)

	payload := map[string]interface{}{
		"username":   "ad",
		"password":   "123",
		"music_path": "/invalid/path",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/init/exec", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, float64(0), resp["code"])
}

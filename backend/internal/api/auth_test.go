package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/config"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/internal/middleware"
	"github.com/homemusic/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthAPITestDB(t *testing.T) {
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
	setupAuthAPITestDB(t)
	gin.SetMode(gin.TestMode)

	user := &model.User{Username: "adminuser", Password: "Password123"}
	err := dao.CreateUser(user)
	assert.NoError(t, err)

	r := gin.Default()
	r.POST("/api/v1/auth/login", Login)

	payload := map[string]interface{}{
		"username": "adminuser",
		"password": "Password123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.NotEmpty(t, resp["data"].(map[string]interface{})["token"])
	assert.Contains(t, w.Header().Get("Set-Cookie"), "auth_token")
}

func TestLoginInvalidCredentials(t *testing.T) {
	setupAuthAPITestDB(t)
	gin.SetMode(gin.TestMode)

	user := &model.User{Username: "adminuser", Password: "Password123"}
	err := dao.CreateUser(user)
	assert.NoError(t, err)

	r := gin.Default()
	r.POST("/api/v1/auth/login", Login)

	payload := map[string]interface{}{
		"username": "adminuser",
		"password": "wrongpass",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(common.CodeUnauthorized), resp["code"])
}

func TestAuthMiddlewareRejectsWithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/api/user/info", middleware.AuthMiddleware(), func(c *gin.Context) {
		common.Success(c, map[string]interface{}{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/api/user/info", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(common.CodeUnauthorized), resp["code"])
}

func TestAuthMiddlewareAllowsWithBearerToken(t *testing.T) {
	setupAuthAPITestDB(t)
	gin.SetMode(gin.TestMode)

	user := &model.User{Username: "adminuser", Password: "Password123"}
	err := dao.CreateUser(user)
	assert.NoError(t, err)

	loginRouter := gin.Default()
	loginRouter.POST("/api/v1/auth/login", Login)
	loginReqBody, _ := json.Marshal(map[string]interface{}{
		"username": "adminuser",
		"password": "Password123",
	})
	loginReq, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginReqBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	loginRouter.ServeHTTP(loginW, loginReq)

	var loginResp map[string]interface{}
	err = json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	token := loginResp["data"].(map[string]interface{})["token"].(string)

	r := gin.Default()
	r.GET("/api/user/info", middleware.AuthMiddleware(), func(c *gin.Context) {
		common.Success(c, map[string]interface{}{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/api/user/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

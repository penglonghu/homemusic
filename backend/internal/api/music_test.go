package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MusicAPITestSuite struct {
	suite.Suite
	db      *gorm.DB
	router  *gin.Engine
}

func (suite *MusicAPITestSuite) SetupTest() {
	// 创建测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// 迁移表
	err = db.AutoMigrate(&model.Music{}, &model.SysInit{})
	assert.NoError(suite.T(), err)

	dao.DB = db
	suite.db = db

	// 设置路由
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	// 注册路由（简化版，用于测试）
	apiGroup := suite.router.Group("/api")
	{
		musicGroup := apiGroup.Group("/music")
		{
			musicGroup.GET("/stats", GetScanStatistics)
			musicGroup.GET("/list", GetMusics)
		}
	}
}

func (suite *MusicAPITestSuite) TestGetScanStatistics() {
	// 创建测试数据
	musicDao := dao.NewMusicDao()
	music := &model.Music{FilePath: "/test.mp3", ScanStatus: "scanned"}
	musicDao.CreateMusic(music)

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/music/stats", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), 200, w.Code)
	assert.Contains(suite.T(), w.Body.String(), `"code":0`)
	assert.Contains(suite.T(), w.Body.String(), `"scanned":1`)
}

func (suite *MusicAPITestSuite) TestGetMusics() {
	// 创建测试数据
	musicDao := dao.NewMusicDao()
	music := &model.Music{FilePath: "/test.mp3", Title: "Test Song"}
	musicDao.CreateMusic(music)

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/music/list?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), 200, w.Code)
	assert.Contains(suite.T(), w.Body.String(), `"code":0`)
	assert.Contains(suite.T(), w.Body.String(), `"Test Song"`)
}

func TestMusicAPITestSuite(t *testing.T) {
	suite.Run(t, new(MusicAPITestSuite))
}
package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MusicServiceTestSuite struct {
	suite.Suite
	db         *gorm.DB
	musicDao   *dao.MusicDao
	musicSvc   *MusicService
	testDir    string
}

func (suite *MusicServiceTestSuite) SetupTest() {
	// 创建测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// 迁移表
	err = db.AutoMigrate(&model.Music{})
	assert.NoError(suite.T(), err)

	// 设置DAO
	dao.DB = db
	suite.db = db
	suite.musicDao = dao.NewMusicDao()
	suite.musicSvc = NewMusicService()

	// 创建临时测试目录
	suite.testDir, err = os.MkdirTemp("", "music_test_*")
	assert.NoError(suite.T(), err)
}

func (suite *MusicServiceTestSuite) TearDownTest() {
	if suite.testDir != "" {
		os.RemoveAll(suite.testDir)
	}
}

func (suite *MusicServiceTestSuite) TestScanMusicDirectory() {
	// 创建测试音频文件
	testFile := filepath.Join(suite.testDir, "test.mp3")
	err := os.WriteFile(testFile, []byte("fake mp3 content"), 0644)
	assert.NoError(suite.T(), err)

	// 执行扫描
	err = suite.musicSvc.ScanMusicDirectory(suite.testDir)
	assert.NoError(suite.T(), err)

	// 验证结果
	music, err := suite.musicDao.GetMusicByPath(testFile)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), music)
	assert.Equal(suite.T(), "scanned", music.ScanStatus)
}

func (suite *MusicServiceTestSuite) TestGetScanStatistics() {
	// 创建测试数据
	music1 := &model.Music{FilePath: "/test1.mp3", ScanStatus: "scanned"}
	music2 := &model.Music{FilePath: "/test2.mp3", ScanStatus: "error"}
	suite.musicDao.CreateMusic(music1)
	suite.musicDao.CreateMusic(music2)

	// 获取统计
	stats, err := suite.musicSvc.GetScanStatistics()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), stats["total"])
	assert.Equal(suite.T(), int64(1), stats["scanned"])
	assert.Equal(suite.T(), int64(1), stats["error"])
}

func (suite *MusicServiceTestSuite) TestGetMusics() {
	// 创建测试数据
	music := &model.Music{FilePath: "/test.mp3", Title: "Test Song", Artist: "Test Artist"}
	suite.musicDao.CreateMusic(music)

	// 获取列表
	musics, total, err := suite.musicSvc.GetMusics(1, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), musics, 1)
	assert.Equal(suite.T(), "Test Song", musics[0].Title)
}

func TestMusicServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MusicServiceTestSuite))
}
package dao

import (
	"github.com/homemusic/backend/config"
	"github.com/homemusic/backend/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
	dsn := config.Conf.Database.Path
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁用表名复数
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	DB = db
	NewInitDao.db = db

	// 自动迁移用户、初始化状态和音乐表
	err = db.AutoMigrate(&model.User{}, &model.SysInit{}, &model.Music{})
	if err != nil {
		return err
	}

	zap.L().Info("database init success")
	return nil
}
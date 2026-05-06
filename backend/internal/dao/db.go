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

	// 设置SQLite连接以支持UTF-8字符编码
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	_, err = sqlDB.Exec("PRAGMA encoding = 'UTF-8';")
	if err != nil {
		zap.L().Error("设置SQLite编码失败", zap.Error(err))
		return err
	}
	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		zap.L().Error("启用外键约束失败", zap.Error(err))
		return err
	}

	// 自动迁移用户、初始化状态、音乐表和播放历史表
	err = db.AutoMigrate(&model.User{}, &model.SysInit{}, &model.Music{}, &model.PlayHistory{})
	if err != nil {
		return err
	}

	zap.L().Info("database init success")
	return nil
}
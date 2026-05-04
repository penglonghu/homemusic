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

	// 🔥 修复：升级 GORM 后，可以直接使用 AutoMigrate，无需再手动判断 HasTable
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	zap.L().Info("database init success")
	return nil
}
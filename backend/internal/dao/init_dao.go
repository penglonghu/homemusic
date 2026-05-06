package dao

import (
	"github.com/homemusic/backend/internal/constant"
	"github.com/homemusic/backend/internal/model"
	"gorm.io/gorm"
)

type InitDao struct {
	db *gorm.DB
}

var NewInitDao = &InitDao{db: DB}

// GetInitStatus 获取系统初始化状态
func (d *InitDao) GetInitStatus() (*model.SysInit, error) {
	var init model.SysInit
	err := d.db.First(&init).Error
	if err == gorm.ErrRecordNotFound {
		return &model.SysInit{IsInit: constant.InitStatusNot}, nil
	}
	return &init, err
}

// CreateInitConfig 创建初始化配置
func (d *InitDao) CreateInitConfig(init *model.SysInit) error {
	return d.db.Create(init).Error
}

// UpdateInitConfig 更新初始化配置
func (d *InitDao) UpdateInitConfig(init *model.SysInit) error {
	return d.db.Save(init).Error
}

// IsInitialized 检查系统初始化状态
func (d *InitDao) IsInitialized() (bool, error) {
	init, err := d.GetInitStatus()
	if err != nil {
		return false, err
	}
	return init.IsInit, nil
}

// SetDB 设置数据库实例，用于测试和运行时注入
func (d *InitDao) SetDB(db *gorm.DB) {
	d.db = db
}

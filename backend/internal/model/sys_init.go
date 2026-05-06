package model

import "time"

// SysInit 系统初始化配置表 兼容SQLite 单主键无冲突
type SysInit struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IsInit    bool      `gorm:"default:false;comment:是否初始化完成" json:"is_init"`
	MusicPath string    `gorm:"type:varchar(512);comment:音乐存储目录" json:"music_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
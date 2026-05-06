package model

import (
	"time"
)

// PlayHistory 用户播放记录
type PlayHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	MusicID   uint      `gorm:"not null;index" json:"music_id"`
	PlayedAt  time.Time `gorm:"not null;index" json:"played_at"`
	Music     Music     `gorm:"foreignKey:MusicID" json:"music"`
}

// TableName 指定表名
func (PlayHistory) TableName() string {
	return "play_history"
}

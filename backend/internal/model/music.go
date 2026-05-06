package model

import (
	"time"

	"gorm.io/gorm"
)

// Music 音乐文件模型
type Music struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 文件信息
	FilePath    string    `gorm:"type:varchar(1024);not null;uniqueIndex" json:"file_path"`
	FileSize    int64     `gorm:"not null" json:"file_size"`
	ModTime     time.Time `gorm:"not null" json:"mod_time"`

	// 元数据
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Artist      string    `gorm:"type:varchar(255);index" json:"artist"`
	Album       string    `gorm:"type:varchar(255);index" json:"album"`
	Genre       string    `gorm:"type:varchar(100)" json:"genre"`
	Year        int       `gorm:"default:0" json:"year"`
	Duration    int       `gorm:"default:0;comment:时长(秒)" json:"duration"`
	TrackNumber int       `gorm:"default:0" json:"track_number"`
	Bitrate     int       `gorm:"default:0;comment:比特率(kbps)" json:"bitrate"`

	// 扫描状态
	ScanStatus  string    `gorm:"type:varchar(20);default:'pending';index;comment:pending/scanned/error" json:"scan_status"`
	LastScanned time.Time `json:"last_scanned"`
	ErrorMsg    string    `gorm:"type:text" json:"error_msg"`
}

// TableName 指定表名
func (Music) TableName() string {
	return "music"
}
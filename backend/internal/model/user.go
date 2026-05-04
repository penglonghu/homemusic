package model

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	// 如果你升级了gorm，保持 `gorm:"primaryKey"` 即可
	// 如果你未升级gorm，请改为 `gorm:"primaryKey;autoIncrement:false"` 绕过Bug
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"size:32;not null;unique"`
	Password  string         `gorm:"size:64;not null"`
	Email     string         `gorm:"size:64;unique"`
	Avatar    string         `gorm:"size:255"`
	Status    int            `gorm:"default:1"`
}

// 密码加密
func (u *User) BeforeSave(*gorm.DB) error {
	if u.Password == "" {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// 密码验证
func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
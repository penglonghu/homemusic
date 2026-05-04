package dao

import (
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetUserByID 根据ID查询用户
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		// 未找到错误不打印错误日志，仅返回
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		common.Logger.Error("查询用户失败", zap.Uint("user_id", id), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名查询用户
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		common.Logger.Error("根据用户名查询用户失败", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func CreateUser(user *model.User) error {
	err := DB.Create(user).Error
	if err != nil {
		common.Logger.Error("创建用户失败", zap.String("username", user.Username), zap.Error(err))
		return err
	}
	common.Logger.Info("创建用户成功", zap.Uint("user_id", user.ID), zap.String("username", user.Username))
	return nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *model.User) error {
	err := DB.Save(user).Error
	if err != nil {
		common.Logger.Error("更新用户失败", zap.Uint("user_id", user.ID), zap.Error(err))
		return err
	}
	common.Logger.Info("更新用户成功", zap.Uint("user_id", user.ID))
	return nil
}

// DeleteUser 软删除用户
func DeleteUser(id uint) error {
	err := DB.Delete(&model.User{}, id).Error
	if err != nil {
		common.Logger.Error("删除用户失败", zap.Uint("user_id", id), zap.Error(err))
		return err
	}
	common.Logger.Info("删除用户成功", zap.Uint("user_id", id))
	return nil
}

// CheckUsernameExist 检查用户名是否存在
func CheckUsernameExist(username string) (bool, error) {
	var count int64
	err := DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		common.Logger.Error("检查用户名存在性失败", zap.String("username", username), zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

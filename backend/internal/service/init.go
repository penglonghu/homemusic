package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/internal/model"
	"go.uber.org/zap"
)

// CheckInitStatus 获取初始化状态
func CheckInitStatus() (map[string]interface{}, error) {
	initStatus, err := dao.NewInitDao.GetInitStatus()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"is_initialized": initStatus.IsInit,
		"music_path":     initStatus.MusicPath,
	}, nil
}

// ExecuteInitialization 执行系统初始化
func ExecuteInitialization(username, password, email, musicPath string) error {
	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("用户名长度必须在3-32个字符之间")
	}
	if len(password) < 6 || len(password) > 64 {
		return fmt.Errorf("密码长度必须在6-64个字符之间")
	}
	if email != "" && len(email) > 64 {
		return fmt.Errorf("邮箱长度不能超过64个字符")
	}
	if musicPath == "" {
		return fmt.Errorf("音乐目录路径不能为空")
	}
	if len(musicPath) > 255 {
		return fmt.Errorf("音乐目录路径过长")
	}
	if !filepath.IsAbs(musicPath) {
		return fmt.Errorf("音乐目录必须使用绝对路径")
	}

	cleanPath := filepath.Clean(musicPath)
	if !strings.HasPrefix(cleanPath, "/") {
		return fmt.Errorf("音乐目录路径不合法")
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("音乐目录不存在")
		}
		return fmt.Errorf("无法访问音乐目录: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("路径不是一个目录")
	}

	// 检查是否已经初始化
	initialized, err := dao.NewInitDao.IsInitialized()
	if err != nil {
		return err
	}
	if initialized {
		return fmt.Errorf("系统已初始化")
	}

	exists, err := dao.CheckUsernameExist(username)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("用户名已存在")
	}

	user := &model.User{
		Username: username,
		Password: password,
		Email:    email,
		Status:   1,
	}

	err = dao.CreateUser(user)
	if err != nil {
		zap.L().Error("创建管理员用户失败", zap.Error(err), zap.String("username", username))
		return fmt.Errorf("创建管理员用户失败")
	}

	initConfig, err := dao.NewInitDao.GetInitStatus()
	if err != nil {
		return err
	}

	if initConfig.ID == 0 {
		initConfig = &model.SysInit{}
	}
	initConfig.IsInit = true
	initConfig.MusicPath = cleanPath
	initConfig.UpdatedAt = time.Now()

	if initConfig.ID == 0 {
		err = dao.NewInitDao.CreateInitConfig(initConfig)
	} else {
		err = dao.NewInitDao.UpdateInitConfig(initConfig)
	}
	if err != nil {
		zap.L().Error("保存初始化状态失败", zap.Error(err), zap.String("music_path", cleanPath))
		return fmt.Errorf("保存初始化状态失败")
	}

	zap.L().Info("系统初始化完成", zap.String("music_path", cleanPath), zap.Uint("admin_user_id", user.ID))
	return nil
}

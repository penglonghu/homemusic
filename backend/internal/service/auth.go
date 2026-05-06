package service

import (
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/homemusic/backend/config"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/pkg/utils"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrIPLocked           = errors.New("当前IP已被锁定，请15分钟后再试")
)

const (
	maxFailedAttempts     = 5
	lockDuration          = 15 * time.Minute
	rememberExpireSeconds = 30 * 24 * 3600
	usernamePattern       = `^[a-zA-Z0-9_]+$`
)

var (
	loginFailures  = make(map[string]*loginFailure)
	failureMutex   sync.RWMutex
	usernameRegexp = regexp.MustCompile(usernamePattern)
)

type loginFailure struct {
	attempts    int
	lockedUntil time.Time
}

type LoginResult struct {
	Token    string `json:"token"`
	Expires  int    `json:"expires"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// Login 执行用户登录
func Login(username, password string, remember bool, clientIP string) (*LoginResult, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if username == "" || !usernameRegexp.MatchString(username) || len(username) < 3 || len(username) > 32 {
		return nil, errors.New("用户名格式不合法")
	}
	if password == "" || len(password) < 6 || len(password) > 64 {
		return nil, errors.New("密码格式不合法")
	}

	if isIPLocked(clientIP) {
		return nil, ErrIPLocked
	}

	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.CheckPassword(password) || user.Status != 1 {
		registerFailedAttempt(clientIP)
		if common.Logger != nil {
			common.Logger.Warn("登录失败", zap.String("username", username), zap.String("client_ip", clientIP))
		}
		return nil, ErrInvalidCredentials
	}

	resetFailedAttempts(clientIP)

	expireSeconds := config.Conf.JWT.Expire
	if remember {
		expireSeconds = rememberExpireSeconds
	}

	token, err := utils.GenerateTokenWithExpire(user.ID, expireSeconds)
	if err != nil {
		if common.Logger != nil {
			common.Logger.Error("生成 JWT Token 失败", zap.Error(err), zap.Uint("user_id", user.ID))
		}
		return nil, err
	}

	if common.Logger != nil {
		common.Logger.Info("用户登录成功", zap.Uint("user_id", user.ID), zap.String("username", username), zap.String("client_ip", clientIP), zap.Bool("remember", remember))
	}

	return &LoginResult{
		Token:    token,
		Expires:  expireSeconds,
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

func isIPLocked(clientIP string) bool {
	failureMutex.RLock()
	defer failureMutex.RUnlock()

	failure, ok := loginFailures[clientIP]
	if !ok {
		return false
	}
	return failure.lockedUntil.After(time.Now())
}

func registerFailedAttempt(clientIP string) {
	failureMutex.Lock()
	defer failureMutex.Unlock()

	failure, ok := loginFailures[clientIP]
	if !ok {
		failure = &loginFailure{}
		loginFailures[clientIP] = failure
	}

	failure.attempts++
	if failure.attempts >= maxFailedAttempts {
		failure.lockedUntil = time.Now().Add(lockDuration)
	}
}

func resetFailedAttempts(clientIP string) {
	failureMutex.Lock()
	defer failureMutex.Unlock()

	delete(loginFailures, clientIP)
}

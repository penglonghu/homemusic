package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/homemusic/backend/config"
)

// CustomClaims 自定义JWT载荷
type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint) (string, error) {
	return GenerateTokenWithExpire(userID, config.Conf.JWT.Expire)
}

// GenerateTokenWithExpire 生成指定过期时间的JWT Token
func GenerateTokenWithExpire(userID uint, expireSeconds int) (string, error) {
	cfg := config.Conf.JWT
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "homemusic",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析JWT Token
func ParseToken(tokenStr string) (*CustomClaims, error) {
	cfg := config.Conf.JWT
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 修复这里！！！原代码写错了
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
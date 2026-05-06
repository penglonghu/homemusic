package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/service"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Remember bool   `json:"remember"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Expires  int    `json:"expires"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// Login 用户登录接口
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if common.Logger != nil {
			common.Logger.Warn("登录参数校验失败", zap.Error(err))
		}
		common.ParamError(c, "请求参数错误")
		return
	}

	result, err := service.Login(req.Username, req.Password, req.Remember, c.ClientIP())
	if err != nil {
		if errors.Is(err, service.ErrIPLocked) {
			if common.Logger != nil {
				common.Logger.Warn("IP 锁定登录请求", zap.String("username", req.Username), zap.String("client_ip", c.ClientIP()))
			}
			common.Fail(c, common.CodeForbidden, err.Error())
			return
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			if common.Logger != nil {
				common.Logger.Warn("用户名或密码错误", zap.String("username", req.Username), zap.String("client_ip", c.ClientIP()))
			}
			common.Fail(c, common.CodeUnauthorized, err.Error())
			return
		}
		if common.Logger != nil {
			common.Logger.Error("登录失败", zap.Error(err))
		}
		common.ServerError(c)
		return
	}

	cookieMaxAge := result.Expires
	c.SetCookie("auth_token", result.Token, cookieMaxAge, "/", "", true, true)

	common.Success(c, LoginResponse{
		Token:    result.Token,
		Expires:  result.Expires,
		UserID:   result.UserID,
		Username: result.Username,
	})
}

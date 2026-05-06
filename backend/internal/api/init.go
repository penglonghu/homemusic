package api

import (
	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/service"
	"go.uber.org/zap"
)

// ExecInitRequest 执行初始化请求
type ExecInitRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=32"`
	Password  string `json:"password" binding:"required,min=6,max=64"`
	Email     string `json:"email" binding:"omitempty,email,max=64"`
	MusicPath string `json:"music_path" binding:"required,max=255"`
}

// CheckInit 初始化状态查询
func CheckInit(c *gin.Context) {
	status, err := service.CheckInitStatus()
	if err != nil {
		common.Logger.Error("获取初始化状态失败", zap.Error(err))
		common.Fail(c, common.CodeServerError, "获取初始化状态失败")
		return
	}
	common.Success(c, status)
}

// ExecInit 执行初始化
func ExecInit(c *gin.Context) {
	var req ExecInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Logger.Warn("初始化参数校验失败", zap.Error(err))
		common.Fail(c, common.CodeParamError, "请求参数错误")
		return
	}

	err := service.ExecuteInitialization(req.Username, req.Password, req.Email, req.MusicPath)
	if err != nil {
		common.Logger.Warn("执行初始化失败", zap.Error(err))
		common.Fail(c, common.CodeParamError, err.Error())
		return
	}

	common.Success(c, map[string]interface{}{
		"message": "初始化完成",
	})
}

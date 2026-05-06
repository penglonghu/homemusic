package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/config"
	"github.com/homemusic/backend/internal/api"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/middleware"
	"go.uber.org/zap"
)

// RunServer 启动服务
func RunServer() {
	cfg := config.Conf.Server

	// 设置gin模式
	gin.SetMode(cfg.Mode)

	// 初始化路由
	r := gin.Default()

	// 注册路由
	apiGroup := r.Group("/api")
	{
		// 公开接口
		apiGroup.GET("/health", api.HealthCheck)

		// 初始化接口
		initGroup := apiGroup.Group("/v1/init")
		{
			initGroup.GET("/check", api.CheckInit)
			initGroup.POST("/exec", api.ExecInit)
		}

		// 需要认证的接口
		authGroup := apiGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			// 用户接口
			authGroup.GET("/user/info", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				common.Success(c, map[string]interface{}{"user_id": userID})
			})

			// 音乐接口
			musicGroup := authGroup.Group("/music")
			{
				musicGroup.POST("/scan", api.ScanMusic)
				musicGroup.GET("/stats", api.GetScanStatistics)
				musicGroup.GET("/list", api.GetMusics)
				musicGroup.POST("/clean", api.CleanDeletedFiles)
			}
		}
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Port)
	common.Logger.Info("server start success", zap.String("addr", addr))
	if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
		common.Logger.Error("server start failed", zap.Error(err))
		panic("server start failed")
	}
}

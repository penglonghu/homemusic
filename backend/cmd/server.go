package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/config"
	"github.com/homemusic/backend/internal/api"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/middleware"
	"github.com/homemusic/backend/pkg/utils"
	"go.uber.org/zap"
)

func RunServer() {
	cfg := config.Conf.Server
	gin.SetMode(cfg.Mode)

	r := gin.Default()
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", api.HealthCheck)

		initGroup := apiGroup.Group("/v1/init")
		{
			initGroup.GET("/check", api.CheckInit)
			initGroup.POST("/exec", api.ExecInit)
		}

		authGroup := apiGroup.Group("/v1/auth")
		{
			authGroup.POST("/login", api.Login)
		}

		authGroup = apiGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.GET("/user/info", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				common.Success(c, map[string]interface{}{"user_id": userID})
			})

			musicGroup := authGroup.Group("/music")
			{
				musicGroup.POST("/scan", api.ScanMusic)
				musicGroup.GET("/stats", api.GetScanStatistics)
				musicGroup.GET("/list", api.GetMusics)
				musicGroup.GET("/browse", api.BrowseMusic)
				musicGroup.GET("/search", api.SearchMusic)
				musicGroup.GET("/recent", api.GetRecentPlayRecords)
				musicGroup.POST("/recent", api.RecordPlayHistory)
				musicGroup.POST("/clean", api.CleanDeletedFiles)
			}
		}

		// 无需认证的公开接口
		publicGroup := apiGroup.Group("")
		{
			publicGroup.GET("/music/stream/:id", api.StreamMusic)  // 音频流接口无需认证
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	certFile := cfg.CertFile
	keyFile := cfg.KeyFile
	if certFile == "" {
		certFile = "./certs/server.crt"
	}
	if keyFile == "" {
		keyFile = "./certs/server.key"
	}

	if err := utils.EnsureSelfSignedCert(certFile, keyFile); err != nil {
		common.Logger.Error("ensure tls cert failed", zap.Error(err))
		panic("ensure tls cert failed")
	}

	common.Logger.Info("https cert ready", zap.String("cert_file", certFile), zap.String("key_file", keyFile))
	common.Logger.Info("server start success", zap.String("addr", addr))
	if err := r.RunTLS(addr, certFile, keyFile); err != nil && err != http.ErrServerClosed {
		common.Logger.Error("server start failed", zap.Error(err))
		panic("server start failed")
	}
}
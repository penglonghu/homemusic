package api

import (
	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/internal/common"
)

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	common.Success(c, map[string]interface{}{
		"status":  "ok",
		"service": "homemusic backend",
	})
}

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/pkg/utils"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先支持 Authorization Header，回退到 auth_token Cookie
		authHeader := c.GetHeader("Authorization")
		var tokenStr string
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else if cookie, err := c.Cookie("auth_token"); err == nil {
			tokenStr = cookie
		}

		if tokenStr == "" {
			common.Unauthorized(c)
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			common.Unauthorized(c)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

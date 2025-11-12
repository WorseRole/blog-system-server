package middleware

import (
	"blog-system-server/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从header中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "未提供认证Token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查 Token格式 “Bearer {token}”
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "Token格式错误",
				"data":    nil,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "无效的Token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(uint64), exists
}

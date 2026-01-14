package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/jwt"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

// Auth 认证中间件
func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证Token")
			c.Abort()
			return
		}

		// 验证Token格式 (Bearer <token>)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 解析Token
		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			switch err {
			case jwt.ErrTokenExpired:
				response.Unauthorized(c, "Token已过期")
			case jwt.ErrTokenNotValidYet:
				response.Unauthorized(c, "Token尚未生效")
			case jwt.ErrTokenMalformed:
				response.Unauthorized(c, "Token格式错误")
			default:
				response.Unauthorized(c, "Token无效")
			}
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（Token存在则验证，不存在则跳过）
func OptionalAuth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := jwtManager.ParseToken(token)
		if err == nil {
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
		}

		c.Next()
	}
}

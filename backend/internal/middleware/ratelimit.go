package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/redis"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 每分钟允许的请求数
	RequestsPerMinute int
	// 是否跳过已认证用户
	SkipAuthenticated bool
}

// RateLimit 限流中间件
// 基于IP地址进行限流，使用Redis存储计数
// 仅针对非登录用户进行限流
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果配置为跳过已认证用户，且用户已认证，则直接放行
		if config.SkipAuthenticated {
			// 检查用户是否已认证（通过检查上下文中是否有userID）
			if _, exists := c.Get("userID"); exists {
				c.Next()
				return
			}
		}

		// 获取客户端IP
		clientIP := c.ClientIP()
		if clientIP == "" {
			response.BadRequest(c, "无法获取客户端IP")
			c.Abort()
			return
		}

		// 构建Redis键
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		// 使用INCR命令原子性增加计数
		ctx := c.Request.Context()
		newCount, err := redis.Get().Incr(ctx, key).Result()
		if err != nil {
			// Redis错误，记录日志但允许请求通过（降级策略）
			c.Next()
			return
		}

		// 如果是新创建的键，设置过期时间
		if newCount == 1 {
			redis.Expire(key, time.Minute)
		}

		// 检查是否超过限制
		if int(newCount) > config.RequestsPerMinute {
			response.TooManyRequests(c, fmt.Sprintf("请求过于频繁，每分钟最多%d次请求", config.RequestsPerMinute))
			c.Abort()
			return
		}

		c.Next()
	}
}
